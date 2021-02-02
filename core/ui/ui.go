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

package ui

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/gorilla/context"
	"infini-framework/core/api"
	"infini-framework/core/api/router"
	"infini-framework/core/global"
	"infini-framework/core/ui/websocket"
	"infini-framework/core/util"
	"net/http"
	_ "net/http/pprof"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var router *httprouter.Router
var mux *http.ServeMux
var l sync.Mutex
var uiConfig *UIConfig

func GetUIConfig() UIConfig {
	return *uiConfig
}

var bindAddress string

func GetBindAddress() string {
	return bindAddress
}

func StartUI(cfg *UIConfig) {
	uiConfig = cfg
	//start web ui
	mux = http.NewServeMux()

	router = httprouter.New(mux)
	//router.RedirectTrailingSlash=false
	//router.RedirectFixedPath=false

	//registered handlers
	if registeredUIHandler != nil {
		for k, v := range registeredUIHandler {
			log.Debug("register custom http handler: ", k)
			mux.Handle(k, v)
		}
	}
	if registeredUIFuncHandler != nil {
		for k, v := range registeredUIFuncHandler {
			log.Debug("register custom http handler: ", k)
			mux.HandleFunc(k, v)
		}
	}
	if registeredUIMethodHandler != nil {
		for k, v := range registeredUIMethodHandler {
			for m, n := range v {
				log.Debug("register custom http handler: ", k, " ", m)
				router.Handle(k, m, n)
			}
		}
	}

	//init websocket,TODO configurable
	websocket.InitWebSocket()
	mux.HandleFunc("/ws", websocket.ServeWs)

	if registeredWebSocketCommandHandler != nil {
		for k, v := range registeredWebSocketCommandHandler {
			log.Debug("register custom websocket handler: ", k, " ", v)
			websocket.HandleWebSocketCommand(k, webSocketCommandUsage[k], v)
		}
	}

	schema := "http://"

	if uiConfig.NetworkConfig.SkipOccupiedPort {
		bindAddress = util.AutoGetAddress(uiConfig.NetworkConfig.GetBindingAddr())
	} else {
		bindAddress = uiConfig.NetworkConfig.GetBindingAddr()
	}

	if uiConfig.TLSConfig.TLSEnabled {
		log.Debug("tls enabled")

		schema = "https://"

		certFile := path.Join(global.Env().SystemConfig.PathConfig.Cert, "*c*rt*")
		match, err := filepath.Glob(certFile)
		if err != nil {
			panic(err)
		}
		if len(match) <= 0 {
			panic(errors.New("no cert file found, the file name must end with .crt"))
		}
		certFile = match[0]

		keyFile := path.Join(global.Env().SystemConfig.PathConfig.Cert, "*key*")
		match, err = filepath.Glob(keyFile)
		if err != nil {
			panic(err)
		}
		if len(match) <= 0 {
			panic(errors.New("no key file found, the file name must end with .key"))
		}
		keyFile = match[0]

		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}

		srv := &http.Server{
			Addr:         bindAddress,
			Handler:      context.ClearHandler(router),
			TLSConfig:    cfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		}

		go func() {
			err = srv.ListenAndServeTLS(certFile, keyFile)
			if err != nil {
				log.Error(err)
				panic(err)
			}
		}()

	} else {
		go func() {
			err := http.ListenAndServe(bindAddress, context.ClearHandler(router))
			if err != nil {
				log.Error(err)
				panic(err)
			}
		}()

	}

	err := util.WaitServerUp(bindAddress, 30*time.Second)
	if err != nil {
		panic(err)
	}

	log.Info("ui server listen at: ", schema, bindAddress)

}

// RegisteredUIHandler is a hub for registered ui handler
var registeredUIHandler map[string]http.Handler

// RegisteredUIFuncHandler is a hub for registered ui handler
var registeredUIFuncHandler map[string]func(http.ResponseWriter, *http.Request)

// RegisteredUIMethodHandler is a hub for registered ui handler
var registeredUIMethodHandler map[string]map[string]func(w http.ResponseWriter, req *http.Request, ps httprouter.Params)

var registeredWebSocketCommandHandler map[string]func(c *websocket.WebsocketConnection, array []string)
var webSocketCommandUsage map[string]string

// HandleUIFunc register ui request handler
func HandleUIFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	l.Lock()
	if registeredUIFuncHandler == nil {
		registeredUIFuncHandler = map[string]func(http.ResponseWriter, *http.Request){}
	}
	registeredUIFuncHandler[pattern] = handler
	l.Unlock()
}

// HandleUI register ui request handler
func HandleUI(pattern string, handler http.Handler) {

	l.Lock()
	if registeredUIHandler == nil {
		registeredUIHandler = map[string]http.Handler{}
	}
	registeredUIHandler[pattern] = handler
	l.Unlock()
}

// HandleUIMethod register ui request handler
func HandleUIMethod(method api.Method, pattern string, handler func(w http.ResponseWriter, req *http.Request, ps httprouter.Params)) {
	l.Lock()
	if registeredUIMethodHandler == nil {
		registeredUIMethodHandler = map[string]map[string]func(w http.ResponseWriter, req *http.Request, ps httprouter.Params){}
	}

	m := registeredUIMethodHandler[string(method)]
	if m == nil {
		registeredUIMethodHandler[string(method)] = map[string]func(w http.ResponseWriter, req *http.Request, ps httprouter.Params){}
	}
	registeredUIMethodHandler[string(method)][pattern] = handler
	l.Unlock()
}

// HandleWebSocketCommand register websocket command handler
func HandleWebSocketCommand(command string, usage string, handler func(c *websocket.WebsocketConnection, array []string)) {

	l.Lock()
	if registeredWebSocketCommandHandler == nil {
		registeredWebSocketCommandHandler = map[string]func(c *websocket.WebsocketConnection, array []string){}
		webSocketCommandUsage = map[string]string{}
	}

	command = strings.ToLower(strings.TrimSpace(command))
	registeredWebSocketCommandHandler[command] = handler
	webSocketCommandUsage[command] = usage
	l.Unlock()
}

// GetPagination return a pagination html code snippet
func GetPagination(from, size, total int, url string, param map[string]interface{}) string {

	//TODO limit when es is the database driver
	//if total > 10000 {
	//	total = 10000
	//}

	if total <= size {
		return ""
	}

	var cur = from / size

	var buffer bytes.Buffer
	buffer.WriteString("<ul class=\"uk-pagination\" data-uk-pagination=\"{items:")
	buffer.WriteString(strconv.Itoa(total))
	buffer.WriteString(", itemsOnPage:")
	buffer.WriteString(strconv.Itoa(size))
	buffer.WriteString(",currentPage:")
	buffer.WriteString(strconv.Itoa(cur))
	buffer.WriteString("}\"></ul>")
	buffer.WriteString("<script type=\"text/javascript\">")

	// init args start
	var moreArgs bytes.Buffer
	moreArgs.WriteString("var args='")
	if len(param) > 0 {
		for k, v := range param {
			hostStr := fmt.Sprintf("&%s=%v", k, v)
			moreArgs.WriteString(hostStr)
		}
	}

	moreArgs.WriteString("';")

	if moreArgs.Len() > 0 {
		buffer.Write(moreArgs.Bytes())
	}

	buffer.WriteString("var size=")
	buffer.WriteString(strconv.Itoa(size))
	buffer.WriteString(";")

	//init args end

	buffer.WriteString("    $(function() {")

	buffer.WriteString("$('[data-uk-pagination]').on('select.uk.pagination', function(e, pageIndex){")

	buffer.WriteString("var from=pageIndex*size;")

	buffer.WriteString("window.location='?from='+from+'&size='+size+args")

	buffer.WriteString("});")

	buffer.WriteString("   });")

	//init para for hot key  start
	buffer.WriteString(fmt.Sprintf("var maxpage = %v;", total))
	if from > 0 && from >= size {
		buffer.WriteString(fmt.Sprintf("var prev_page='?from=%v&size='+size+args;", from-size))

	}
	if from+size < total {
		buffer.WriteString(fmt.Sprintf("var next_page='?from=%v&size='+size+args;", from+size))
	}
	//init para for hot key end

	buffer.WriteString("</script>")

	return buffer.String()
}
