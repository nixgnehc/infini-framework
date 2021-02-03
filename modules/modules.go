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

package modules

import (
	"github.com/nixgnehc/infini-framework/core/module"
	"github.com/nixgnehc/infini-framework/modules/api"
	"github.com/nixgnehc/infini-framework/modules/boltdb"
	"github.com/nixgnehc/infini-framework/modules/cluster"
	"github.com/nixgnehc/infini-framework/modules/elastic"
	"github.com/nixgnehc/infini-framework/modules/filter"
	"github.com/nixgnehc/infini-framework/modules/pipeline"
	"github.com/nixgnehc/infini-framework/modules/queue"
	"github.com/nixgnehc/infini-framework/modules/stats"
	"github.com/nixgnehc/infini-framework/modules/ui"
)

// RegisterSystemModule is where modules are registered
func Register() {
	//module.RegisterSystemModule(nsq.NSQModule{})
	module.RegisterSystemModule(elastic.ElasticModule{})
	module.RegisterSystemModule(boltdb.StorageModule{})
	module.RegisterSystemModule(filter.FilterModule{})
	module.RegisterSystemModule(stats.SimpleStatsModule{})
	module.RegisterSystemModule(queue.DiskQueue{})
	module.RegisterSystemModule(api.APIModule{})
	module.RegisterSystemModule(ui.UIModule{})
	module.RegisterSystemModule(pipeline.PipeModule{})
	module.RegisterSystemModule(cluster.ClusterModule{})
}
