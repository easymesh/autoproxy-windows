package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func AutoRunningGet() bool {
	if DataIntValueGet("autorunning") > 0 {
		return true
	}
	return false
}

func AutoRunningSet(flag bool) {
	if flag {
		DataIntValueSet("autorunning", 1)
	} else {
		DataIntValueSet("autorunning", 0)
	}
}

func BaseSetting() {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var auto *walk.RadioButton

	_, err := Dialog{
		AssignTo:      &dlg,
		Title:         "Base Setting",
		Icon:          walk.IconShield(),
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		Size:          Size{250, 200},
		MinSize:       Size{250, 200},
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "Auto Startup: ",
					},
					RadioButton{
						AssignTo: &auto,
						OnBoundsChanged: func() {
							auto.SetChecked(AutoRunningGet())
						},
						OnClicked: func() {
							auto.SetChecked(!AutoRunningGet())
							AutoRunningSet(!AutoRunningGet())
						},
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							go func() {
								dlg.Accept()
							}()
						},
					},
					PushButton{
						AssignTo: &cancelPB,
						Text:     "Cancel",
						OnClicked: func() {
							dlg.Cancel()
						},
					},
				},
			},
		},
	}.Run(mainWindow)

	if err != nil {
		logs.Error(err.Error())
	}
}
