package statsd

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/quipo/statsd"
	. "infini-framework/core/config"
	"infini-framework/core/env"
	"infini-framework/core/errors"
	"infini-framework/core/stats"
	"sync"
	"time"
)

type StatsDConfig struct {
	Host              string        `config:"host"`
	Port              int           `config:"port"`
	Namespace         string        `config:"namespace"`
	IntervalInSeconds time.Duration `config:"interval_in_seconds"`
}
type StatsDModule struct {
}

var cfg *Config

func (module StatsDModule) Setup(config *Config) {
	cfg = config
}

var statsdInited bool
var statsdclient *statsd.StatsdClient
var buffer *statsd.StatsdBuffer
var l1 sync.RWMutex

var defaultStatsdConfig = StatsDConfig{
	Host:              "localhost",
	Port:              8125,
	Namespace:         "app.",
	IntervalInSeconds: 1,
}

func (module StatsDModule) Name() string {
	return "statsd"
}

func (module StatsDModule) Start() error {
	if statsdInited {
		panic(errors.New("statsd not inited"))
	}

	config := defaultStatsdConfig
	//cfg.Unpack(&config)
	env.ParseConfig("statsd", &config)

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	l1.Lock()
	defer l1.Unlock()
	statsdclient = statsd.NewStatsdClient(addr, config.Namespace)

	log.Debug("statsd connec to, ", addr, ",prefix:", config.Namespace)

	err := statsdclient.CreateSocket()
	if nil != err {
		log.Warn(err)
		return err
	}

	interval := time.Second * config.IntervalInSeconds // aggregate stats and flush every 2 seconds
	buffer = statsd.NewStatsdBuffer(interval, statsdclient)

	statsdInited = true

	stats.Register(module)
	return nil
}

func (module StatsDModule) Stop() error {
	if statsdclient != nil {
		statsdclient.Close()
	}
	return nil
}

func (module StatsDModule) Absolute(category, key string, value int64) {

	if !statsdInited {
		return
	}
	buffer.Absolute(category+"."+key, value)
}

func (module StatsDModule) Increment(category, key string) {

	module.IncrementBy(category, key, 1)
}

func (module StatsDModule) IncrementBy(category, key string, value int64) {
	if !statsdInited {
		return
	}
	buffer.Incr(category+"."+key, value)
}

func (module StatsDModule) Decrement(category, key string) {
	module.DecrementBy(category, key, 1)
}

func (module StatsDModule) DecrementBy(category, key string, value int64) {
	if !statsdInited {
		return
	}
	buffer.Decr(category+"."+key, value)
}

func (module StatsDModule) Timing(category, key string, v int64) {
	if !statsdInited {
		return
	}
	buffer.Timing(category+"."+key, v)

}

func (module StatsDModule) Gauge(category, key string, v int64) {
	if !statsdInited {
		return
	}
	buffer.Gauge(category+"."+key, v)
}

func (module StatsDModule) Stat(category, key string) int64 {
	return 0
}

func (module StatsDModule) StatsAll() *[]byte {
	return nil
}
