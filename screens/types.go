package screens

import (
	"yoru/screens/components"
	"yoru/types"
)

type manager struct {
	types.ScreenManager
	tabBar *components.TabBar
}

type root struct {
	types.Screen
	tabs   []types.Screen
	tabBar *components.TabBar
}

type home struct {
	types.Screen
}
