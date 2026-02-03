package screens

import (
	"yoru/models"
	"yoru/screens/components"
	"yoru/screens/forms"
	"yoru/screens/popups"
	"yoru/terminal"
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

type terminalScreen struct {
	types.Screen
	hostID          uint
	host            *models.Host
	emulator        *terminal.Emulator
	connectionPopup *popups.ConnectionPopup
	connecting      bool
	connected       bool
	connectionLog   *models.ConnectionLog
	keyCaptureMode  types.KeyCaptureMode
}

type focusArea int
