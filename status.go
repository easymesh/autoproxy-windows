package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var statusFlow *walk.StatusBarItem
var sessionFlow *walk.StatusBarItem
var requestFlow *walk.StatusBarItem

func UpdateStatFlow(speed string, session string, request string) {
	if statusFlow != nil {
		statusFlow.SetText("Flow: " + speed)
		requestFlow.SetText("Request: " + request)
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
		{
			AssignTo: &requestFlow,
			Text:     "",
			Width:    80,
		},
	}
}
