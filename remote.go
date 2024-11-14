package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/easymesh/autoproxy-windows/engin"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func remoteOptions() []string {
	var output []string
	list := ConfigGet().RemoteList
	for _, v := range list {
		output = append(output, v.Name)
	}
	if len(output) == 0 {
		output = append(output, "")
	}
	return output
}

func remoteIndex() int {
	remote := ConfigGet().RemoteName
	list := ConfigGet().RemoteList
	for i, v := range list {
		if v.Name == remote {
			return i
		}
	}
	return 0
}

func remoteFind(name string) Remote {
	list := ConfigGet().RemoteList
	for _, v := range list {
		if v.Name == name {
			return v
		}
	}
	return Remote{
		Name: name, Protocol: "HTTPS",
	}
}

func remoteDelete(name string) {
	remoteList := ConfigGet().RemoteList
	for i, v := range remoteList {
		if v.Name == name {
			remoteList = append(remoteList[:i], remoteList[i+1:]...)
		}
	}
	RemoteListSave(remoteList)
}

func remoteListUpdate(item Remote) {
	remoteList := ConfigGet().RemoteList

	for i, v := range remoteList {
		if v.Name == item.Name {
			remoteList[i] = item
			RemoteListSave(remoteList)
			return
		}
	}
	remoteList = append(remoteList, item)
	RemoteListSave(remoteList)
}

func protocolOptions() []string {
	return []string{
		"HTTP", "HTTPS",
	}
}

func protocolIndex(protocol string) int {
	if protocol == "HTTP" {
		return 0
	} else {
		return 1
	}
}

func TestEngin(testhttps string, item *Remote) (time.Duration, error) {
	now := time.Now()
	if !engin.IsConnect(item.Address, 5) {
		return 0, fmt.Errorf("remote address connnect %s fail", item.Address)
	}

	urls, err := url.Parse(testhttps)
	if err != nil {
		logs.Error("%s raw url parse fail, %s", testhttps, err.Error())
		return 0, err
	}

	var auth *engin.AuthInfo
	if item.Auth {
		auth = &engin.AuthInfo{
			User:  item.User,
			Token: item.Password,
		}
	}

	var tls bool
	if strings.ToLower(item.Protocol) == "https" {
		tls = true
	}

	forward, err := engin.NewHttpsProtocol(item.Address, 10, auth, tls, "", "")
	if err != nil {
		logs.Error("new remote http proxy fail, %s", err.Error())
		return 0, err
	}

	defer forward.Close()

	request, err := http.NewRequest("GET", testhttps, nil)
	if err != nil {
		logs.Error("%s raw url parse fail, %s", testhttps, err.Error())
		return 0, err
	}

	if strings.ToLower(urls.Scheme) == "https" {
		_, err = forward.Https(engin.Address(urls), request)
	} else {
		_, err = forward.Http(request)
	}

	if err != nil && err.Error() != "EOF" {
		logs.Error("remote server %s forward %s fail, %s",
			item.Address, urls.RawPath, err.Error())
		return 0, err
	}

	return time.Since(now), nil
}

func RemoteServer() {
	var remoteDlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var remote, protocol *walk.ComboBox
	var auth *walk.RadioButton
	var user, passwd, address, testurl *walk.LineEdit
	var testButton *walk.PushButton
	var remoteCurrent Remote

	remoteList := ConfigGet().RemoteList
	if len(remoteList) > 0 {
		remoteCurrent = remoteList[0]
	}

	var updateHandler = func() {
		protocol.SetCurrentIndex(protocolIndex(remoteCurrent.Protocol))
		address.SetText(remoteCurrent.Address)
		auth.SetChecked(remoteCurrent.Auth)
		user.SetEnabled(remoteCurrent.Auth)
		passwd.SetEnabled(remoteCurrent.Auth)
		user.SetText(remoteCurrent.User)
		passwd.SetText(remoteCurrent.Password)
	}

	_, err := Dialog{
		AssignTo:      &remoteDlg,
		Title:         "Remote Proxy Setting",
		Icon:          walk.IconShield(),
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		Size:          Size{250, 300},
		MinSize:       Size{250, 300},
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "Remote Proxy: ",
					},
					ComboBox{
						AssignTo:     &remote,
						Editable:     true,
						CurrentIndex: 0,
						Model:        remoteOptions(),
						OnBoundsChanged: func() {
							remoteCurrent = remoteFind(remote.Text())
							updateHandler()
						},
						OnCurrentIndexChanged: func() {
							remoteCurrent = remoteFind(remote.Text())
							updateHandler()
						},
						OnEditingFinished: func() {
							remoteCurrent = remoteFind(remote.Text())
							updateHandler()
						},
					},

					Label{
						Text: "Remote Address: ",
					},

					LineEdit{
						AssignTo:    &address,
						Text:        remoteCurrent.Address,
						ToolTipText: "192.168.1.3:8080",
						OnEditingFinished: func() {
							remoteCurrent.Address = address.Text()
						},
					},

					Label{
						Text: "Remote protocol: ",
					},
					ComboBox{
						AssignTo:     &protocol,
						Model:        protocolOptions(),
						CurrentIndex: protocolIndex(remoteCurrent.Protocol),
						OnCurrentIndexChanged: func() {
							remoteCurrent.Protocol = protocol.Text()
						},
					},
					Label{
						Text: "Auth: ",
					},
					RadioButton{
						AssignTo: &auth,
						OnBoundsChanged: func() {
							auth.SetChecked(remoteCurrent.Auth)
						},
						OnClicked: func() {
							auth.SetChecked(!remoteCurrent.Auth)
							remoteCurrent.Auth = !remoteCurrent.Auth

							user.SetEnabled(remoteCurrent.Auth)
							passwd.SetEnabled(remoteCurrent.Auth)
						},
					},
					Label{
						Text: "User: ",
					},
					LineEdit{
						AssignTo: &user,
						Text:     remoteCurrent.User,
						Enabled:  remoteCurrent.Auth,
						OnEditingFinished: func() {
							remoteCurrent.User = user.Text()
						},
					},
					Label{
						Text: "Password: ",
					},
					LineEdit{
						AssignTo: &passwd,
						Text:     remoteCurrent.Password,
						Enabled:  remoteCurrent.Auth,
						OnEditingFinished: func() {
							remoteCurrent.Password = passwd.Text()
						},
					},
					PushButton{
						AssignTo: &testButton,
						Text:     "Testing",
						OnClicked: func() {
							go func() {
								testButton.SetEnabled(false)
								delay, err := TestEngin(testurl.Text(), &remoteCurrent)
								if err != nil {
									ErrorBoxAction(remoteDlg, err.Error())
								} else {
									info := fmt.Sprintf("%s, %s %dms",
										"Test Pass",
										"Delay", delay/time.Millisecond)
									InfoBoxAction(remoteDlg, info)
								}
								testButton.SetEnabled(true)
							}()
						},
					},
					LineEdit{
						AssignTo: &testurl,
						Text:     ConfigGet().TestUrl,
						OnEditingFinished: func() {
							TestUrlSave(testurl.Text())
						},
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &acceptPB,
						Text:     "Save",
						OnClicked: func() {
							if remoteCurrent.Auth {
								if remoteCurrent.User == "" || remoteCurrent.Password == "" {
									ErrorBoxAction(remoteDlg, "Please input user and passwd")
									return
								}
							}
							if remoteCurrent.Name == "" || remoteCurrent.Address == "" {
								ErrorBoxAction(remoteDlg, "Please input name and address")
								return
							}
							remoteListUpdate(remoteCurrent)
							remoteDlg.Accept()
							ConsoleRemoteUpdate()
						},
					},
					PushButton{
						AssignTo: &cancelPB,
						Text:     "Delete",
						OnClicked: func() {
							remoteDelete(remoteCurrent.Name)

							remoteList := ConfigGet().RemoteList
							if len(remoteList) > 0 {
								remoteCurrent = remoteList[0]
							}

							remote.SetModel(remoteOptions())
							remote.SetCurrentIndex(0)
							updateHandler()

							ConsoleRemoteUpdate()
						},
					},
					PushButton{
						AssignTo: &cancelPB,
						Text:     "Cancel",
						OnClicked: func() {
							remoteDlg.Cancel()
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
