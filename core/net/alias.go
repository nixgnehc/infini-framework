package net

import (
	log "github.com/cihub/seelog"
	"infini-framework/core/errors"
	"infini-framework/core/util"
)

func checkPermission() {
	log.Debug("to continue use net alias, you need to run as root or elevate with sudo.")
	if !util.RequireSudo() {
		panic(errors.New("root or sudo permission needed."))
	}
}
