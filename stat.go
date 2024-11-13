package main

import (
	"fmt"

	"github.com/easymesh/autoproxy-windows/engin"
	"github.com/lxn/walk"
)

func StatUpdate(stat engin.StatInfo) {
	UpdateStatFlow(
		ByteView(stat.ForwardSize), fmt.Sprintf("%d", stat.SessionCnt), fmt.Sprintf("%d", stat.RequestCnt))
	NotifyUpdateFlow(ByteView(stat.ForwardSize))
}

func StatRunningStatus(enable bool) {
	var image *walk.Icon
	if enable {
		image = ICON_Network_Enable
	} else {
		image = ICON_Network_Disable
	}
	UpdateStatFlag(image)
	NotifyUpdateIcon(image)
}
