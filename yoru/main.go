package main

import (
	"yoru/screens"
	"yoru/types"
	"yoru/utils/bridge"
	"yoru/utils/errors"
)

func main() {
	if err := bridge.New(screens.ScreenManager, types.WithAltScreen); err != nil {
		errors.ExitOnBridgeFailedStart(err)
	}
}
