package screens

import "yoru/types"

type manager struct {
	types.ScreenManager
	Current types.Screen
}

type root struct {
	types.Screen
	toggled bool
}
