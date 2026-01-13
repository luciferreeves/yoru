package screens

import (
	"yoru/types"
)

type manager struct {
	types.ScreenManager
	tabBar types.TabBar
}

type home struct {
	types.Screen
	navBar types.NavBar
}
