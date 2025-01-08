package main

import (
	"net"

	"github.com/astaxie/beego/logs"
)

func ModeOptions() []string {
	return []string{
		"Full Local Proxy",
		"Match Domain Proxy",
		"Auto Connect Proxy",
		"Full Remote Proxy",
	}
}

func ModeOptionsIdx() int {
	mode := ConfigGet().Mode
	if int(mode) < len(ModeOptions()) {
		return int(mode)
	}
	return 0
}

func ModeOptionsSet(idx int) {
	ModeSave(ModeType(idx))
}

func IfaceOptions() []string {
	output := []string{"0.0.0.0", "::"}
	ifaces, err := net.Interfaces()
	if err != nil {
		logs.Error("get interfaces fail, %s", err.Error())
		return output
	}
	for _, v := range ifaces {
		if v.Flags&net.FlagUp == 0 {
			continue
		}
		address, err := InterfaceAddsGet(&v)
		if err != nil {
			continue
		}
		for _, addr := range address {
			output = append(output, addr.String())
		}
	}
	return output
}
