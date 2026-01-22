package components

import (
	"yoru/screens/styles"
	"yoru/shared"
	"yoru/types"

	"github.com/charmbracelet/lipgloss"
)

var TabBar = &tabBar{}

func (tabBar *tabBar) AddTab(tab types.Tab) {
	tabBar.tabs = append(tabBar.tabs, tab)
	if len(tabBar.tabs) == 1 {
		tabBar.activeIndex = 0
	}
}

func (tabBar *tabBar) GetCurrentScreen() types.Screen {
	if tabBar.activeIndex < 0 || tabBar.activeIndex >= len(tabBar.tabs) {
		return nil
	}

	return tabBar.tabs[tabBar.activeIndex].Screen
}

func (tabBar *tabBar) UpdateCurrentScreen(screen types.Screen) {
	if tabBar.activeIndex >= 0 && tabBar.activeIndex < len(tabBar.tabs) {
		tabBar.tabs[tabBar.activeIndex].Screen = screen
	}
}

func (tabBar *tabBar) SwitchToTab(index int) {
	if index < 0 || index >= len(tabBar.tabs) {
		return
	}

	tabBar.activeIndex = index
}

func (tabBar *tabBar) NextTab() {
	totalTabs := len(tabBar.tabs)
	if totalTabs > 1 {
		tabBar.activeIndex = (tabBar.activeIndex + 1) % totalTabs
	}
}

func (tabBar *tabBar) PrevTab() {
	totalTabs := len(tabBar.tabs)
	if totalTabs > 1 {
		tabBar.activeIndex = tabBar.activeIndex - 1
		if tabBar.activeIndex < 0 {
			tabBar.activeIndex = totalTabs - 1
		}
	}
}

func (tabBar *tabBar) Render() string {
	if len(tabBar.tabs) == 0 {
		return ""
	}

	measureWidth := lipgloss.Width
	var renderedTabs []string

	for index, tab := range tabBar.tabs {
		if index == tabBar.activeIndex {
			content := styles.ActiveTab.Render(" " + tab.Name + " ")
			renderedTabs = append(renderedTabs, content)
		} else {
			content := styles.InactiveTab.Render(" " + tab.Name + " ")
			renderedTabs = append(renderedTabs, content)
		}
	}

	tabsContent := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	remainingWidth := shared.GlobalState.ScreenWidth - measureWidth(tabsContent)

	if remainingWidth > 0 {
		gap := styles.TabBarBackground.Width(remainingWidth).Render("")
		tabsContent = lipgloss.JoinHorizontal(lipgloss.Top, tabsContent, gap)
	}

	return styles.TabBarBackground.Width(shared.GlobalState.ScreenWidth).Render(tabsContent)
}
