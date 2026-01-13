package components

import (
	"yoru/screens/styles"
	"yoru/shared"

	"github.com/charmbracelet/lipgloss"
)

var NavBar = &navBar{
	items:       []string{"Hosts", "SFTP", "Keychain", "Known Hosts", "Logs", "Preferences"},
	activeIndex: 0,
}

func (navBar *navBar) GetActiveTab() (int, string) {
	return navBar.activeIndex, navBar.items[navBar.activeIndex]
}

func (navBar *navBar) SwitchToTab(index int) {
	if index < 0 || index >= len(navBar.items) {
		return
	}

	navBar.activeIndex = index
}

func (navBar *navBar) NextTab() {
	totalItems := len(navBar.items)
	navBar.activeIndex = (navBar.activeIndex + 1) % totalItems
}

func (navBar *navBar) PrevTab() {
	totalItems := len(navBar.items)
	navBar.activeIndex--
	if navBar.activeIndex < 0 {
		navBar.activeIndex = totalItems - 1
	}
}

func (navBar *navBar) Render() string {
	measureWidth := lipgloss.Width
	var renderedItems []string

	for index, item := range navBar.items {
		if index == navBar.activeIndex {
			renderedItems = append(renderedItems, styles.ActiveNavBar.Render(item))
		} else {
			renderedItems = append(renderedItems, styles.InactiveNavBar.Render(item))
		}
	}

	navContent := lipgloss.JoinHorizontal(lipgloss.Top, renderedItems...)
	remainingWidth := shared.GlobalState.ScreenWidth - measureWidth(navContent)

	if remainingWidth > 0 {
		gap := styles.NavBarBackground.Width(remainingWidth).Render("")
		navContent = lipgloss.JoinHorizontal(lipgloss.Top, navContent, gap)
	}

	return styles.NavBarBackground.Width(shared.GlobalState.ScreenWidth).Render(navContent)
}
