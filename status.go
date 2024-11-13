package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var statusFlow *walk.StatusBarItem
var connectFlow *walk.StatusBarItem

func UpdateStatFlow(flow string) {
	if statusFlow != nil {
		statusFlow.SetText(flow)
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
			AssignTo:    &statusFlow,
			Icon:        ICON_Network_Disable,
			ToolTipText: "status",
			Width:       80,
		},
		{
			AssignTo:    &connectFlow,
			Text:        "0",
			ToolTipText: "sessions",
			Width:       80,
		},
	}
}
