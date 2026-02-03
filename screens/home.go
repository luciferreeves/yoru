package screens

import (
	"yoru/screens/components"
	"yoru/screens/styles"
	"yoru/shared"
	"yoru/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var homeScreen = &home{
	navBar: components.NavBar,
}

func (screen *home) Init() tea.Cmd {
	hostsScreen.Init()
	keychainScreen.Init()
	logsScreen.Init()
	return nil
}

func (screen *home) Update(msg tea.Msg) (types.Screen, tea.Cmd) {
	switch message := msg.(type) {
	case tea.KeyMsg:
		if cmd := screen.OnKeyPress(message); cmd != nil {
			return screen, cmd
		}
	}

	currentIndex, _ := screen.navBar.GetActiveTab()
	switch currentIndex {
	case 0:
		_, cmd := hostsScreen.Update(msg)
		return screen, cmd
	case 2:
		_, cmd := keychainScreen.Update(msg)
		return screen, cmd
	case 4:
		_, cmd := logsScreen.Update(msg)
		return screen, cmd
	}

	return screen, nil
}

func (screen *home) View() string {
	navBarView := screen.navBar.Render()

	var contentText string
	currentIndex, _ := screen.navBar.GetActiveTab()
	switch currentIndex {
	case 0:
		contentText = hostsScreen.View()
	case 1:
		contentText = "SFTP content coming soon"
	case 2:
		contentText = keychainScreen.View()
	case 3:
		contentText = "Known Hosts content coming soon"
	case 4:
		contentText = logsScreen.View()
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

func (screen *home) OnKeyPress(key tea.KeyMsg) tea.Cmd {
	currentIndex, _ := screen.navBar.GetActiveTab()

	if currentIndex == 0 && (key.Type == tea.KeyLeft || key.Type == tea.KeyRight) {
		if hostsScreen.focusedArea == formFocus || hostsScreen.sidebar.IsFilterActive() || hostsScreen.deletePopup.IsVisible() {
			return nil
		}
	}

	if currentIndex == 2 && (key.Type == tea.KeyLeft || key.Type == tea.KeyRight) {
		if keychainScreen.focusedArea == keychainFormFocus || keychainScreen.sidebar.IsFilterActive() || keychainScreen.deletePopup.IsVisible() {
			return nil
		}
	}

	switch key.Type {
	case tea.KeyLeft:
		screen.navBar.PrevTab()
	case tea.KeyRight:
		screen.navBar.NextTab()
	}

	return nil
}
