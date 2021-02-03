package net

import (
	log "github.com/cihub/seelog"
	"github.com/nixgnehc/infini-framework/core/errors"
	"github.com/nixgnehc/infini-framework/core/util"
)

func checkPermission() {
	log.Debug("to continue use net alias, you need to run as root or elevate with sudo.")
	if !util.RequireSudo() {
		panic(errors.New("root or sudo permission needed."))
	}
}
