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

func ConsoleIndex() int {
	interfaces := IfaceOptions()
	listenAddress := ConfigGet().Address
	for i, v := range interfaces {
		if v == listenAddress {
			return i
		}
	}
	return 0
}

func ConsoleEnable(enable bool) {
	consoleIface.SetEnabled(enable)
	consolePort.SetEnabled(enable)
}

func ConsoleRemoteUpdate() {
	remote := ConfigGet().RemoteName
	remoteList := ConfigGet().RemoteList

	var remoteOptions []string
	consoleRemoteProxy.SetCurrentIndex(0)
	for i, v := range remoteList {
		if v.Name == remote {
			consoleRemoteProxy.SetCurrentIndex(i)
		}
		remoteOptions = append(remoteOptions, v.Name)
	}
	consoleRemoteProxy.SetModel(remoteOptions)
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

	go func() {
		for {
			if active != nil && active.Visible() {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		if len(ConfigGet().RemoteList) > 0 || ConfigGet().Mode == "local" {
			activeFunc()
		}
	}()

	return []Widget{
		Label{
			Text: "Listen Address: ",
		},
		ComboBox{
			AssignTo:     &consoleIface,
			CurrentIndex: ConsoleIndex(),
			Model:        IfaceOptions(),
			OnCurrentIndexChanged: func() {
				ListenAddressSave(consoleIface.Text())
			},
			OnBoundsChanged: func() {
				consoleIface.SetCurrentIndex(ConsoleIndex())
			},
		},
		Label{
			Text: "Listen Port: ",
		},
		NumberEdit{
			AssignTo:    &consolePort,
			Value:       float64(ConfigGet().Port),
			ToolTipText: "1~65535",
			MaxValue:    65535,
			MinValue:    1,
			OnValueChanged: func() {
				ListenPortSave(int(consolePort.Value()))
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
			CurrentIndex: remoteIndex(),
			OnBoundsChanged: func() {
				if len(ConfigGet().RemoteList) == 0 {
					consoleMode.SetCurrentIndex(0)
					ModeOptionsSet(0)
					consoleMode.SetEnabled(false)
				} else {
					consoleMode.SetEnabled(true)
				}
			},
			OnCurrentIndexChanged: func() {
				if len(ConfigGet().RemoteList) == 0 {
					consoleMode.SetCurrentIndex(0)
					ModeOptionsSet(0)
					consoleMode.SetEnabled(false)
				} else {
					consoleMode.SetEnabled(true)
				}

				consoleRemoteProxy.SetEnabled(false)
				RemoteSave(consoleRemoteProxy.Text())

				go func() {
					err := RemoteForwardUpdate()
					if err != nil {
						ErrorBoxAction(mainWindow, err.Error())
					}
					consoleRemoteProxy.SetEnabled(true)
				}()
			},
			Model: remoteOptions(),
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
