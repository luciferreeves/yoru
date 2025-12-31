package main

import (
	"yoru/screens"
	"yoru/types"
	"yoru/utils/bridge"
)

func main() {
	if err := bridge.New(screens.ScreenManager, types.WithAltScreen); err != nil {
		bridge.ExitOnError(err)
	}
}
