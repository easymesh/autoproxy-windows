package main

import (
	"sync"

	"github.com/easymesh/autoproxy-windows/engin"
)

type LocalAccessInfo struct {
	Hostname string
	Access   bool
}

type AutoCtrl struct {
	sync.RWMutex
	cache map[string]LocalAccessInfo
}

var autoCtrl AutoCtrl

func init() {
	autoCtrl.cache = make(map[string]LocalAccessInfo, 100)
}

func AutoCheck(address string) bool {
	autoCtrl.RLock()
	result, ok := autoCtrl.cache[address]
	autoCtrl.RUnlock()
	if ok {
		return result.Access
	}

	result.Hostname = address
	result.Access = engin.IsConnect(address, 3)

	autoCtrl.Lock()
	autoCtrl.cache[address] = result
	autoCtrl.Unlock()

	return result.Access
}

func AutoCheckUpdate(address string, access bool) {
	autoCtrl.RLock()
	result, ok := autoCtrl.cache[address]
	autoCtrl.RUnlock()

	if !ok || result.Access != access {
		autoCtrl.Lock()
		autoCtrl.cache[address] = LocalAccessInfo{address, access}
		autoCtrl.Unlock()
	}
}
