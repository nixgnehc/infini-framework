/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pipeline

import (
	log "github.com/cihub/seelog"
	"github.com/nixgnehc/infini-framework/core/errors"
	"github.com/nixgnehc/infini-framework/core/global"
	"github.com/nixgnehc/infini-framework/core/stats"
	"github.com/nixgnehc/infini-framework/core/util"
	"runtime"
	"strings"
	"time"
)

type Pipeline struct {
	id string

	name string

	input   Input
	output  Output
	filters []Filter

	startProcessor Processor

	processors []Processor

	runningState RunningState

	context *Context

	endProcessor Processor

	errorProcessor Processor

	currentProcessor string

	newVersion bool
}

func NewPipeline(name string) *Pipeline {
	pipe := &Pipeline{}
	pipe.id = util.GetUUID()
	pipe.filters = []Filter{}
	pipe.name = strings.TrimSpace(name)
	pipe.context = &Context{}
	pipe.context.init()
	return pipe
}

func (pipe *Pipeline) Context(s *Context) *Pipeline {
	if s != nil {
		pipe.context = s
		pipe.context.init()
	}

	return pipe
}

func (pipe *Pipeline) GetID() string {
	return pipe.id
}

func (pipe *Pipeline) GetContext() *Context {
	return pipe.context
}

func (pipe *Pipeline) Input(s Input) *Pipeline {
	pipe.input = s
	pipe.newVersion=true
	return pipe
}

func (pipe *Pipeline) Output(s Output) *Pipeline {
	pipe.output = s
	return pipe
}

func (pipe *Pipeline) Filter(s Filter) *Pipeline {
	pipe.filters = append(pipe.filters, s)
	return pipe
}

func (pipe *Pipeline) Start(s Processor) *Pipeline {
	pipe.startProcessor = s
	pipe.processors = []Processor{}
	return pipe
}

func (pipe *Pipeline) Join(s Processor) *Pipeline {
	pipe.processors = append(pipe.processors, s)
	return pipe
}

func (pipe *Pipeline) Error(s Processor) *Pipeline {
	pipe.errorProcessor = s
	return pipe
}

func (pipe *Pipeline) End(s Processor) *Pipeline {
	pipe.endProcessor = s
	return pipe
}

// setCurrentProcessor set current joint's name, used for debugging
func (context *Pipeline) setCurrentProcessor(name string) {
	context.currentProcessor = name
}

func (pipe *Pipeline) CurrentProcessor() string {
	return pipe.currentProcessor
}

func (pipe *Pipeline) Start1() *Pipeline {
	pipe.runningState = STARTED
	return pipe
}

func (pipe *Pipeline) Pause() *Pipeline {
	pipe.runningState = PAUSED
	return pipe
}

func (pipe *Pipeline) Resume() *Pipeline {
	pipe.runningState = STARTED
	return pipe
}

func (pipe *Pipeline) Stop() *Pipeline {
	pipe.runningState = STOPPED
	return pipe
}

func (pipe *Pipeline) Run1() *Pipeline {

	var err error
	if pipe.output != nil {
		//open output first
		err = pipe.output.Open()
		if err != nil {
			panic(err)
		}
	}

	if pipe.input == nil {
		panic(errors.New("pipeline input can't be null"))
	}

	//open input after output
	err = pipe.input.Open()
	if err != nil {
		panic(err)
	}

	pipe.Start1()

	var data []byte

	for {
		switch pipe.runningState {
		case STARTED:
			data, err = pipe.input.Read()
			if err != nil {
				panic(err)
			}

			for _, f := range pipe.filters {
				err = f.Filter(data)
				if err != nil {
					panic(err)
				}
			}

			if pipe.output != nil {
				err = pipe.output.Write(data)
				if err != nil {
					panic(err)
				}
			}
			pipe.runningState = FINISHED

			return nil
		case PAUSED:
			break
		case STOPPED:
			err = pipe.input.Close()
			if err != nil {
				panic(err)
			}

			if pipe.output != nil {
				err = pipe.output.Close()
				if err != nil {
					panic(err)
				}
			}

			return nil
		}

	}
}

func (pipe *Pipeline) Run() *Context {

	if pipe.newVersion{
		pipe.Run1()
		return nil
	}

	pipe.runningState = STARTED

	stats.Increment(pipe.name+".pipeline", "total")

	//final phrase
	defer func() {
		if !global.Env().IsDebug {
			if r := recover(); r != nil {

				if r == nil {
					return
				}
				var v string
				switch r.(type) {
				case error:
					v = r.(error).Error()
				case runtime.Error:
					v = r.(runtime.Error).Error()
				case string:
					v = r.(string)
				}
				//pipe.context.Set(CONTEXT_TASK_Status, TaskInterrupted)
				//pipe.context.Set(CONTEXT_TASK_Message, util.ToJson(v, false))

				log.Error("error in pipeline, ", pipe.name, ", ", pipe.id, ", ", pipe.currentProcessor, ", ", v)
				stats.Increment(pipe.name+".pipeline", "error")
			}
		}

		if !pipe.context.IsExit() && (!(pipe.context.IgnoreBroken && pipe.context.IsEnd())) {
			pipe.endPipeline()
		}

		stats.Increment(pipe.name+".pipeline", "finished")
	}()

	var err error

	pipe.startPipeline()

	for _, v := range pipe.processors {
		log.Trace("pipe, ", pipe.name, ", start joint,", v.Name())
		if pipe.context.IsEnd() {
			log.Trace("break joint,", v.Name())
			stats.Increment(pipe.name+".pipeline", "break")
			return pipe.context
		}

		if pipe.context.IsExit() {
			if global.Env().IsDebug {
				log.Debug(util.ToJson(pipe.id, true))
				log.Debug(util.ToJson(pipe.name, true))
				log.Debug(util.ToJson(pipe.context, true))
			}
			log.Trace("exit joint,", v.Name())
			stats.Increment(pipe.name+".pipeline", "exit")
			return pipe.context
		}

		pipe.setCurrentProcessor(v.Name())
		startTime := time.Now().UTC()
		err = v.Process(pipe.context)

		elapsedTime := time.Now().UTC().Sub(startTime)
		stats.Timing(pipe.name+".pipeline", v.Name(), elapsedTime.Nanoseconds())
		if err != nil {
			stats.Increment(pipe.name+".pipeline", "error")
			log.Debugf("%s-%s: %v", pipe.name, v.Name(), err)
			pipe.context.Payload = err.Error()
			pipe.handlePipelineError()
			return pipe.context
		}
		log.Trace(pipe.name, ", end joint,", v.Name())
	}

	return pipe.context
}

func (pipe *Pipeline) startPipeline() {

	log.Trace("start pipeline: ", pipe.name)
	if pipe.startProcessor != nil {
		pipe.setCurrentProcessor(pipe.startProcessor.Name())
		pipe.startProcessor.Process(pipe.context)
	}
	log.Trace("pipeline: ", pipe.name, ", started")
}

func (pipe *Pipeline) endPipeline() {
	if pipe.context.IsExit() {
		log.Debug("exit pipeline, ", pipe.name, ", ", pipe.context.Payload)
		return
	}

	log.Trace("start finish pipeline, ", pipe.name)
	if pipe.endProcessor != nil {
		pipe.setCurrentProcessor(pipe.endProcessor.Name())
		pipe.endProcessor.Process(pipe.context)
	}
	log.Trace("end finish pipeline, ", pipe.name)
}

func (pipe *Pipeline) handlePipelineError() {

	if pipe.errorProcessor != nil {
		log.Trace("start handle pipeline error, ", pipe.name)
		pipe.setCurrentProcessor(pipe.errorProcessor.Name())
		pipe.errorProcessor.Process(pipe.context)
		log.Trace("end handle pipeline error, ", pipe.name)
	}
}

func NewPipelineFromConfig(name string, config *PipelineConfig, context *Context) *Pipeline {
	if global.Env().IsDebug {
		log.Debugf("pipeline config: %v", util.ToJson(config, true))
	}

	pipe := NewPipeline(name)

	pipe.Context(context)

	if config.Input != nil && config.Input.Enabled {
		input := GetInputJointInstance(config.Input)
		pipe.Input(input)
	}

	if config.Output != nil && config.Output.Enabled {
		input := GetOutputJointInstance(config.Output)
		pipe.Output(input)
	}

	for _, cfg := range config.Filters {
		if cfg.Enabled {
			j := GetFilterJointInstance(cfg)
			pipe.Filter(j)
		}
	}

	if config.StartProcessor != nil && config.StartProcessor.Enabled {
		input := GetJointInstance(config.StartProcessor)
		pipe.Start(input)
	}

	if config.ErrorProcessor != nil && config.ErrorProcessor.Enabled {
		input := GetJointInstance(config.ErrorProcessor)
		pipe.Error(input)
	}

	for _, cfg := range config.Processors {
		if cfg.Enabled {
			j := GetJointInstance(cfg)
			pipe.Join(j)
		}
	}

	if config.EndProcessor != nil && config.EndProcessor.Enabled {
		output := GetJointInstance(config.EndProcessor)
		pipe.End(output)
	}

	if global.Env().IsDebug {
		log.Debugf("get pipeline: %v", util.ToJson(pipe, true))
	}

	return pipe
}
