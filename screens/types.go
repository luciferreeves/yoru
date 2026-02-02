package screens

import (
	"yoru/models"
	"yoru/screens/components"
	"yoru/screens/forms"
	"yoru/screens/popups"
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
	sidebar              *components.HostsSidebar
	form                 *forms.HostForm
	focusedArea          focusArea
	filterWasActive      bool
	deletePopup          *popups.DeleteHostPopup
	identityChooserPopup *popups.IdentityChooserPopup
}

type logs struct {
	types.Screen
	logs        []models.ConnectionLog
	selectedIdx int
}

type focusArea int
