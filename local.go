package main

import (
	"net"

	"github.com/astaxie/beego/logs"
)

type Options struct {
	Name   string
	Detail string
}

var modeOptions = []*Options{
	{"local", "Local Forward"},
	{"auto", "Auto Forward"},
	{"proxy", "Global Forward"},
}

func ModeOptions() []*Options {
	return modeOptions
}

func ModeOptionGet() string {
	return ConfigGet().Mode
}

func ModeOptionsIdx() int {
	mode := ConfigGet().Mode
	for i, opt := range modeOptions {
		if opt.Name == mode {
			return i
		}
	}
	return 0
}

func ModeOptionsSet(idx int) {
	for i, opt := range modeOptions {
		if idx == i {
			ModeSave(opt.Name)
			return
		}
	}
}

func IfaceOptions() []string {
	output := []string{"0.0.0.0", "::"}
	ifaces, err := net.Interfaces()
	if err != nil {
		logs.Error(err.Error())
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
