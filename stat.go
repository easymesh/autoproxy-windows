package main

import (
	"fmt"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var defaultFont = Font{
	PointSize: 14,
	Bold:      true,
}

var LastUpdate time.Time

func StatUpdate(requst, flowsize uint64) {
	flowsizeStr := fmt.Sprintf("%s/s",
		ByteView(int64(float64(flowsize))))

	UpdateStatFlow(flowsizeStr)
	NotifyUpdateFlow(flowsizeStr)
}

func StatInit() error {
	return nil
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
