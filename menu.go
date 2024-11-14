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
					Text: "Domain Setting",
					OnTriggered: func() {
						RemodeEdit()
					},
				},
				Action{
					Text: "Proxy Setting",
					OnTriggered: func() {
						RemoteServer()
					},
				},
				Action{
					Text: "Runlog",
					OnTriggered: func() {
						OpenBrowserWeb(RunlogDirGet())
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
