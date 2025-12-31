package shared

import (
	"yoru/types"
	"yoru/utils/term"
)

var GlobalState types.GlobalState

func init() {
	width, height := term.GetTermSize()
	GlobalState = types.GlobalState{
		ScreenWidth:  width,
		ScreenHeight: height,
		BuildDate:    Date,
		BuildVersion: Version,
	}
}
