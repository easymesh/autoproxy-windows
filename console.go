package main

import (
	"sync"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var consoleIface *walk.ComboBox
var consoleRemoteProxy *walk.ComboBox
var consoleMode *walk.ComboBox
var consolePort *walk.NumberEdit

func ConsoleEnable(enable bool) {
	consoleIface.SetEnabled(enable)
	consolePort.SetEnabled(enable)
}

func ConsoleRemoteUpdate() {
	consoleRemoteProxy.SetModel(RemoteOptions())
	consoleRemoteProxy.SetCurrentIndex(RemoteIndexGet())
}

func ConsoleWidget() []Widget {
	var active *walk.PushButton

	mutex := new(sync.Mutex)

	activeFunc := func() {
		mutex.Lock()

		if ServerRunning() {
			err := ServerShutdown()
			if err != nil {
				ErrorBoxAction(mainWindow, err.Error())
			}
			ConsoleRemoteUpdate()
		} else {
			err := ServerStart()
			if err != nil {
				ErrorBoxAction(mainWindow, err.Error())
			}
		}

		if ServerRunning() {
			StatRunningStatus(true)
			ConsoleEnable(false)
			active.SetImage(ICON_Stop)
			active.SetText("Stop")
		} else {
			StatRunningStatus(false)
			ConsoleEnable(true)
			active.SetImage(ICON_Start)
			active.SetText("Start")
		}

		mutex.Unlock()
	}

	if AutoRunningGet() {
		go func() {
			for {
				if active != nil && active.Visible() {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
			activeFunc()
		}()
	}

	return []Widget{
		Label{
			Text: "Listen Address: ",
		},
		ComboBox{
			AssignTo:     &consoleIface,
			CurrentIndex: LocalIfaceOptionsIdx(),
			Model:        IfaceOptions(),
			OnCurrentIndexChanged: func() {
				LocalIfaceOptionsSet(consoleIface.Text())
			},
		},
		Label{
			Text: "Listen Port: ",
		},
		NumberEdit{
			AssignTo:    &consolePort,
			Value:       float64(PortOptionGet()),
			ToolTipText: "1~65535",
			MaxValue:    65535,
			MinValue:    1,
			OnValueChanged: func() {
				PortOptionSet(int(consolePort.Value()))
			},
		},
		Label{
			Text: "Proxy Mode: ",
		},
		ComboBox{
			AssignTo:      &consoleMode,
			BindingMember: "Name",
			DisplayMember: "Detail",
			CurrentIndex:  ModeOptionsIdx(),
			Model:         ModeOptions(),
			OnCurrentIndexChanged: func() {
				ModeOptionsSet(consoleMode.CurrentIndex())
				go func() {
					ModeUpdate()
				}()
			},
		},
		Label{
			Text: "Remote Proxy: ",
		},
		ComboBox{
			AssignTo:     &consoleRemoteProxy,
			CurrentIndex: RemoteIndexGet(),
			OnBoundsChanged: func() {
				if len(RemoteList()) == 0 {
					consoleMode.SetCurrentIndex(0)
					ModeOptionsSet(0)
					consoleMode.SetEnabled(false)
				} else {
					consoleMode.SetEnabled(true)
				}
			},
			OnCurrentIndexChanged: func() {
				if len(RemoteList()) == 0 {
					consoleMode.SetCurrentIndex(0)
					ModeOptionsSet(0)
					consoleMode.SetEnabled(false)
				} else {
					consoleMode.SetEnabled(true)
				}

				consoleRemoteProxy.SetEnabled(false)
				RemoteIndexSet(consoleRemoteProxy.Text())
				go func() {
					err := RemoteForwardUpdate()
					if err != nil {
						ErrorBoxAction(mainWindow, err.Error())
					}
					consoleRemoteProxy.SetEnabled(true)
				}()
			},
			Model: RemoteOptions(),
		},
		VSpacer{},
		PushButton{
			AssignTo: &active,
			Image:    ICON_Start,
			Text:     "Start",
			OnClicked: func() {
				go activeFunc()
			},
		},
	}
}
