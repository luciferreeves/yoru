package screens

import (
	"yoru/screens/components"
	"yoru/screens/forms"
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

type hosts struct {
	types.Screen
	sidebar         *components.HostsSidebar
	form            *forms.HostForm
	focusedArea     focusArea
	filterWasActive bool
}

type focusArea int
