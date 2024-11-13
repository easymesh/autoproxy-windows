package main

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var mainWindow *walk.MainWindow

var mainWindowWidth = 400
var mainWindowHeight = 180

func waitWindows() {
	for {
		if mainWindow != nil && mainWindow.Visible() {
			mainWindow.SetSize(walk.Size{
				mainWindowWidth,
				mainWindowHeight})
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	NotifyInit()
}

func MainWindowsClose() {
	if mainWindow != nil {
		mainWindow.Close()
		mainWindow = nil
	}
}

func statusUpdate() {
	StatUpdate(StatGet())
}

func init() {
	go func() {
		waitWindows()
		for {
			statusUpdate()
			time.Sleep(time.Second)
		}
	}()
}

var isAuth *walk.RadioButton
var protocal *walk.RadioButton

func mainWindows() {
	CapSignal(CloseWindows)
	cnt, err := MainWindow{
		Title:          "Auto Proxy " + VersionGet(),
		Icon:           ICON_Main,
		AssignTo:       &mainWindow,
		MinSize:        Size{mainWindowWidth, mainWindowHeight - 1},
		Size:           Size{mainWindowWidth, mainWindowHeight - 1},
		Layout:         VBox{Margins: Margins{Top: 10, Bottom: 10, Left: 10, Right: 10}},
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
