package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/astaxie/beego/logs"
	"github.com/easymesh/autoproxy-windows/engin"
)

var access engin.Access

func StatGet() engin.StatInfo {
	acc := access
	if acc != nil {
		return acc.StatGet()
	}
	return engin.StatInfo{}
}

var LocalForward engin.Forward
var RemoteForward engin.Forward

var mutex sync.Mutex

func LocalForwardFunc(address string, r *http.Request) engin.Forward {
	return LocalForward
}

func ProxyForwardFunc(address string, r *http.Request) engin.Forward {
	return RemoteForward
}

func DomainForwardFunc(address string, r *http.Request) engin.Forward {
	if RouteCheck(address) {
		logs.Info("match %s domain and forward to remote proxy", address)
		return RemoteForward
	}
	return LocalForward
}

func AutoForwardFunc(address string, r *http.Request) engin.Forward {
	if AutoCheck(address) {
		logs.Info("auto connect check %s and forward to remote proxy", address)
		return RemoteForward
	}
	return LocalForward
}

func RemoteForwardUpdate() error {
	mutex.Lock()
	defer mutex.Unlock()

	return remoteUpdate()
}

func remoteUpdate() error {
	remoteList := ConfigGet().RemoteList
	if len(remoteList) == 0 {
		return nil
	}

	var remoteCurrent Remote
	var find bool
	for _, item := range remoteList {
		if item.Name == ConfigGet().RemoteName {
			remoteCurrent = item
			find = true
		}
	}

	if !find {
		return fmt.Errorf("the remote proxy server does not exist")
	}

	logs.Info("ready switch to the remote server %v", remoteCurrent)

	var tlsEnable bool
	if strings.ToLower(remoteCurrent.Protocol) == "https" {
		tlsEnable = true
	}

	var auth *engin.AuthInfo
	if remoteCurrent.Auth {
		auth = &engin.AuthInfo{User: remoteCurrent.User, Token: remoteCurrent.Password}
	}

	forward, err := engin.NewHttpsProtocol(remoteCurrent.Address, 60, auth, tlsEnable, "", "")
	if err != nil {
		logs.Error("new remote http proxy fail, %s", err.Error())
		return err
	}

	logs.Info("switch to the remote server %s success", remoteCurrent.Name)

	if RemoteForward != nil {
		RemoteForward.Close()
	}

	RemoteForward = forward
	return nil
}

func modeUpdate() error {
	acc := access
	if acc == nil {
		logs.Warn("server has been stop, mode update disable")
		return nil
	}

	mode := ConfigGet().Mode

	logs.Info("mode switch to %s", mode)

	switch mode {
	case ModeLocal:
		acc.ForwardHandlerSet(LocalForwardFunc)
	case ModeDomain:
		acc.ForwardHandlerSet(DomainForwardFunc)
	case ModeAuto:
		acc.ForwardHandlerSet(AutoForwardFunc)
	case ModeProxy:
		acc.ForwardHandlerSet(ProxyForwardFunc)
	default:
		return fmt.Errorf("mode %d not support", mode)
	}

	logs.Info("server mode switch to %s success", mode)
	return nil
}

func ModeUpdate() error {
	mutex.Lock()
	defer mutex.Unlock()

	return modeUpdate()
}

func ServerStart() error {
	mutex.Lock()
	defer mutex.Unlock()

	var err error

	if access != nil {
		logs.Error("server has been start")
		return fmt.Errorf("server has been start")
	}

	address := ConfigGet().Address
	if !strings.Contains(address, ":") {
		address = fmt.Sprintf("%s:%d", address, ConfigGet().Port)
	} else {
		address = fmt.Sprintf("[%s]:%d", address, ConfigGet().Port)
	}

	logs.Info("server start %s", address)

	access, err = engin.NewHttpsAccess(address, 60, false, "", "")
	if err != nil {
		logs.Error(err.Error())
		return err
	}

	LocalForward = engin.NewDefault(60)

	err = remoteUpdate()
	if err != nil {
		return err
	}

	modeUpdate()

	logs.Info("server start %s success", address)
	return nil
}

func ServerRunning() bool {
	mutex.Lock()
	defer mutex.Unlock()

	if access == nil {
		return false
	}
	return true
}

func ServerShutdown() error {
	mutex.Lock()
	defer mutex.Unlock()

	if access == nil {
		return fmt.Errorf("server has been stop")
	}
	err := access.Shutdown()
	if err != nil {
		logs.Error("shutdown fail, %s", err.Error())
		return err
	}
	access = nil

	LocalForward.Close()
	LocalForward = nil
	return nil
}
