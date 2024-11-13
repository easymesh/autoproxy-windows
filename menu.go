package main

import (
	. "github.com/lxn/walk/declarative"
)

func MenuBarInit() []MenuItem {
	return []MenuItem{
		Menu{
			Text: "Setting",
			Items: []MenuItem{
				Action{
					Text: "Auto Startup",
					OnTriggered: func() {
						BaseSetting()
					},
				},
				Action{
					Text: "Runlog",
					OnTriggered: func() {
						OpenBrowserWeb(logDirGet())
					},
				},
				Separator{},
				Action{
					Text: "Exit",
					OnTriggered: func() {
						CloseWindows()
					},
				},
			},
		},
		Action{
			Text: "Domain",
			OnTriggered: func() {
				RemodeEdit()
			},
		},
		Action{
			Text: "Proxy",
			OnTriggered: func() {
				RemoteServer()
			},
		},
		Action{
			Text: "Mini Windows",
			OnTriggered: func() {
				Notify()
			},
		},
		Action{
			Text: "Sponsor",
			OnTriggered: func() {
				AboutAction()
			},
		},
	}
}
