package main

import (
	rice "github.com/GeertJohan/go.rice"
	"github.com/astaxie/beego/logs"
)

var box *rice.Box

func BoxInit() {
	var err error
	conf := rice.Config{
		LocateOrder: []rice.LocateMethod{rice.LocateEmbedded},
	}
	box, err = conf.FindBox("static")
	if err != nil {
		logs.Error("box init fail, %s", err.Error())
	}
}

func BoxFile() *rice.Box {
	return box
}
