package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var statusFlow *walk.StatusBarItem
var sessionFlow *walk.StatusBarItem

func UpdateStatFlow(speed string, session string) {
	if statusFlow != nil {
		statusFlow.SetText("FlowByte: " + speed)
		sessionFlow.SetText("Session: " + session)
	}
}

func UpdateStatFlag(image *walk.Icon) {
	if statusFlow != nil {
		statusFlow.SetIcon(image)
	}
}

func StatusBarInit() []StatusBarItem {
	return []StatusBarItem{
		{
			AssignTo: &statusFlow,
			Icon:     ICON_Network_Disable,
			Text:     "",
			Width:    100,
		},
		{
			AssignTo: &sessionFlow,
			Text:     "",
			Width:    80,
		},
	}
}
