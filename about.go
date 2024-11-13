package main

import (
	"os/exec"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func OpenBrowserWeb(url string) {
	cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	err := cmd.Run()
	if err != nil {
		logs.Error("run cmd fail, %s", err.Error())
	}
}

var image1 walk.Image
var image2 walk.Image

func LoadImage(name string) walk.Image {
	body, err := BoxFile().Bytes(name)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	file := DEFAULT_HOME + "\\" + name
	err = SaveToFile(file, body)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	image, err := walk.NewImageFromFile(file)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	return image
}

func AboutAction() {
	var ok *walk.PushButton
	var about *walk.Dialog
	var err error

	if image1 == nil {
		image1 = LoadImage("sponsor1.jpg")
	}

	if image2 == nil {
		image2 = LoadImage("sponsor2.jpg")
	}

	_, err = Dialog{
		AssignTo:      &about,
		Title:         "Sponsor",
		Icon:          walk.IconInformation(),
		MinSize:       Size{Width: 300, Height: 200},
		DefaultButton: &ok,
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{
						MinSize: Size{Width: 10},
					},
					ImageView{
						ToolTipText: "Alipay",
						Image:       image1,
						MaxSize:     Size{80, 80},
					},
					HSpacer{
						MinSize: Size{Width: 10},
					},
					ImageView{
						ToolTipText: "WecartPay",
						Image:       image2,
						MaxSize:     Size{80, 80},
					},
					HSpacer{
						MinSize: Size{Width: 10},
					},
				},
			},
			PushButton{
				Text: "paypal.me",
				OnClicked: func() {
					OpenBrowserWeb("https://paypal.me/lixiangyun")
				},
			},
			PushButton{
				Text:      "OK",
				OnClicked: func() { about.Cancel() },
			},
		},
	}.Run(mainWindow)

	if err != nil {
		logs.Error(err.Error())
	}
}
