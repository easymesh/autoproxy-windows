package main

import (
	"fmt"

	"github.com/easymesh/autoproxy-windows/engin"
	"github.com/lxn/walk"
)

func StatUpdate(stat engin.StatInfo) {
	UpdateStatFlow(ByteView(stat.Size), fmt.Sprintf("%d", stat.Session))
	NotifyUpdateFlow(ByteView(stat.Size))
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
