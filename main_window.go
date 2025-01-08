package main

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var mainWindow *walk.MainWindow

func MainWindowsClose() {
	if mainWindow != nil {
		mainWindow.Close()
		mainWindow = nil
	}
}

func init() {
	go func() {
		for {
			StatUpdate(StatGet())
			time.Sleep(time.Second)
		}
	}()
}

func mainWindows() {
	CapSignal(CloseWindows)

	logs.Info("main windows ready to startup")

	cnt, err := MainWindow{
		Title:          "Auto Proxy " + VersionGet(),
		Icon:           ICON_Main,
		AssignTo:       &mainWindow,
		MinSize:        Size{Width: 300, Height: 150},
		Size:           Size{Width: 300, Height: 150},
		Layout:         VBox{Margins: Margins{Top: 5, Bottom: 5, Left: 5, Right: 5}},
		MenuItems:      MenuBarInit(),
		StatusBarItems: StatusBarInit(),
		Children: []Widget{
			Composite{
				Layout:   Grid{Columns: 2},
				Children: ConsoleWidget(),
			},
		},
	}.Run()

	if err != nil {
		logs.Error(err.Error())
	} else {
		logs.Info("main windows exit %d", cnt)
	}

	if err := recover(); err != nil {
		logs.Error(err)
	}

	CloseWindows()
}

func CloseWindows() {
	if ServerRunning() {
		ServerShutdown()
	}
	MainWindowsClose()
	NotifyExit()
}
