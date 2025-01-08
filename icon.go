package main

import (
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
)

func IconLoadFromBox(filename string, size walk.Size) *walk.Icon {
	body, err := BoxFile().Bytes(filename)
	if err != nil {
		logs.Error(err.Error())
		return walk.IconApplication()
	}
	dir := filepath.Join(DEFAULT_HOME, "icon")
	_, err = os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 644)
		if err != nil {
			logs.Error("mkdir %s fail, %s", dir, err.Error())
			return walk.IconApplication()
		}
	}
	err = SaveToFile(filepath.Join(dir, filename), body)
	if err != nil {
		logs.Error("save to file fail, %s", err.Error())
		return walk.IconApplication()
	}
	icon, err := walk.NewIconFromFileWithSize(filepath.Join(dir, filename), size)
	if err != nil {
		logs.Error("new icon from file fail, %s", err.Error())
		return walk.IconApplication()
	}
	return icon
}

var ICON_Main *walk.Icon
var ICON_Network_Disable *walk.Icon
var ICON_Network_Enable *walk.Icon
var ICON_Start *walk.Icon
var ICON_Stop *walk.Icon

var ICON_Max_Size = walk.Size{
	Width: 128, Height: 128,
}

var ICON_Min_Size = walk.Size{
	Width: 16, Height: 16,
}

func IconInit() {
	ICON_Main = IconLoadFromBox("main.ico", ICON_Max_Size)
	ICON_Network_Disable = IconLoadFromBox("network_disable.ico", ICON_Min_Size)
	ICON_Network_Enable = IconLoadFromBox("network_enable.ico", ICON_Min_Size)
	ICON_Start = IconLoadFromBox("start.ico", ICON_Min_Size)
	ICON_Stop = IconLoadFromBox("stop.ico", ICON_Min_Size)
}
