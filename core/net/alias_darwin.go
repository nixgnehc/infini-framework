package net

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/nixgnehc/infini-framework/core/errors"
	"github.com/nixgnehc/infini-framework/core/global"
	"github.com/nixgnehc/infini-framework/core/util"
	"os/exec"
)

//OSX
//sudo /sbin/ifconfig en0 alias 192.168.3.213  255.255.255.0
//sudo /sbin/ifconfig en0 -alias 192.168.3.213

//Linux
//sudo /sbin/ifconfig eth0:1 192.168.3.198 netmask 255.255.255.0
//sudo /sbin/ifconfig eth0:1 down

//Windows
//netsh -c Interface ip add address name="INFINI Ethernet" addr=10.10.0.21 mask=255.255.0.0
//netsh -c Interface ip delete address name="INFINI Ethernet" addr=10.10.0.21
func SetupAlias(device, ip, netmask string) error {

	EnableAlias(device, ip, netmask)

	//register global callback to disable alias before shutdown
	global.RegisterShutdownCallback(func() {
		DisableAlias(device, ip, netmask)
	})

	return nil
}

func EnableAlias(device, ip string, netmask string) error {

	checkPermission()

	log.Debugf("setup net alias %s, %s, %s", device, ip, netmask)
	setupVIP := exec.Command("/usr/bin/sudo", "/sbin/ifconfig", device, "alias", ip, netmask)
	_, err := setupVIP.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to set alias on interface: %s", err))
	}

	ok, err := util.CheckIPBinding(ip)
	if !ok || err != nil {
		return errors.New(fmt.Sprintf("failed to get interface by alias %s: %s\n", device, err))
	}

	log.Debug("net alias was successfully setup/enabled")

	return nil
}

//OSX
//sudo /sbin/ifconfig en0 -alias 192.168.3.213

//Linux
///sbin/ifdown device

//Windows
//netsh interface set interface name="INFINI Ethernet" admin=DISABLED
func DisableAlias(device, ip string, netmask string) error {
	checkPermission()

	setupVIP := exec.Command("/usr/bin/sudo", "/sbin/ifconfig", device, "-alias", ip)
	_, err := setupVIP.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to disable alias on interface: %s", err))
	}

	ok, err := util.CheckIPBinding(ip)
	if ok || err != nil {
		return errors.New(fmt.Sprintf("failed to disable alias on interface %s: %s\n", device, err))
	}

	log.Debug("net alias was successfully disabled")

	return nil
}
