/*
Copyright Medcl (m AT medcl.net)

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

package leader

//import (
//	"infini-framework/core/rpc"
//	."infini-framework/core/cluster/pb"
//	"context"
//	log "github.com/cihub/seelog"
//)
//
//type HealthCheck struct {
//
//}
//
//func (c *HealthCheck) Ping(ctx context.Context, in *HealthCheckRequest) (*HealthCheckResponse, error) {
//
//	log.Info(in.NodeIp,",",in.NodeName,",",in.NodePort)
//
//	return &HealthCheckResponse{Success:true}, nil
//}
//
//func Init() {
//	mys := &HealthCheck{}
//	RegisterHealthCheckServer(rpc.GetRPCServer(), mys)
//
//}