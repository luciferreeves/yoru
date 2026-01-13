package screens

import (
	"yoru/screens/components"
	"yoru/screens/styles"
	"yoru/shared"
	"yoru/types"

	"github.com/charmbracelet/lipgloss"
)

var homeScreen = &home{
	navBar: components.NavBar,
}

func (screen *home) Init() types.Command {
	return nil
}

func (screen *home) Update(event types.Event) (types.Screen, types.Command) {
	switch message := event.(type) {
	case types.KeyPress:
		return screen, screen.OnKeyPress(message)
	}
	return screen, nil
}

func (screen *home) View() string {
	navBarView := screen.navBar.Render()

	var contentText string
	currentIndex, _ := screen.navBar.GetActiveTab()
	switch currentIndex {
	case 0:
		contentText = "Hosts content coming soon"
	case 1:
		contentText = "SFTP content coming soon"
	case 2:
		contentText = "Keychain content coming soon"
	case 3:
		contentText = "Known Hosts content coming soon"
	case 4:
		contentText = "Logs content coming soon"
	case 5:
		contentText = "Preferences content coming soon"
	}

	navBarHeight := lipgloss.Height(navBarView)
	contentHeight := shared.GlobalState.ScreenHeight - navBarHeight - 1 - 2
	contentWidth := shared.GlobalState.ScreenWidth - 2

	contentArea := styles.ContentArea.
		Width(contentWidth).
		Height(contentHeight).
		Render(contentText)

	return lipgloss.JoinVertical(lipgloss.Left, navBarView, contentArea)
}

func (screen *home) OnKeyPress(key types.KeyPress) types.Command {
	switch key.Type {
	case types.KeyLeft:
		screen.navBar.PrevTab()
	case types.KeyRight:
		screen.navBar.NextTab()
	}

	return nil
}
