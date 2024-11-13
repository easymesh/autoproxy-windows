package main

import (
	"net"

	"github.com/astaxie/beego/logs"
)

type Options struct {
	Name   string
	Detail string
}

func ModeOptions() []*Options {
	return []*Options{
		{"local", "Local Forward"},
		{"auto", "Auto Forward"},
		{"proxy", "Global Forward"},
	}
}

func ModeOptionGet() string {
	return ModeOptions()[ModeOptionsIdx()].Name
}

func ModeOptionsIdx() int {
	return int(DataIntValueGet("LocalMode"))
}

func ModeOptionsSet(idx int) {
	err := DataIntValueSet("LocalMode", uint32(idx))
	if err != nil {
		logs.Error(err.Error())
	}
}

func PortOptionGet() int {
	value := DataIntValueGet("LocalPort")
	if value == 0 {
		value = 8080
	}
	return int(value)
}

func PortOptionSet(value int) {
	err := DataIntValueSet("LocalPort", uint32(value))
	if err != nil {
		logs.Error(err.Error())
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

func LocalIfaceOptionsIdx() int {
	ifaces := IfaceOptions()
	ifaceName := DataStringValueGet("LocalIface")
	for idx, v := range ifaces {
		if v == ifaceName {
			return idx
		}
	}
	return 0
}

func LocalIfaceOptionsSet(ifaceName string) {
	err := DataStringValueSet("LocalIface", ifaceName)
	if err != nil {
		logs.Error(err.Error())
	}
}
