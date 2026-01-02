package screens

import (
	"yoru/screens/components"
	"yoru/shared"
	"yoru/types"

	"github.com/charmbracelet/lipgloss"
)

var ScreenManager = &manager{
	tabBar: &components.TabBar{},
}

func (manager *manager) Init() types.Command {
	manager.tabBar.AddTab(types.Tab{
		Name:   "Home",
		Screen: homeScreen,
	})

	return homeScreen.Init()
}

func (manager *manager) Update(event types.Event) (types.Screen, types.Command) {
	switch message := event.(type) {
	case types.KeyPress:
		if command := manager.OnKeyPress(message); command != nil {
			return manager, command
		}
	case types.WindowResized:
		shared.GlobalState.ScreenWidth = message.Width
		shared.GlobalState.ScreenHeight = message.Height
		return manager, nil
	}

	activeScreen := manager.tabBar.ActiveScreen()
	if activeScreen != nil {
		updatedScreen, command := activeScreen.Update(event)
		manager.updateActiveScreen(updatedScreen)
		return manager, command
	}

	return manager, nil
}

func (manager *manager) View() string {
	activeScreen := manager.tabBar.ActiveScreen()
	tabBarView := manager.tabBar.Render()

	contentHeight := shared.GlobalState.ScreenHeight - 1

	var contentView string
	if activeScreen != nil {
		contentView = activeScreen.View()
	}

	contentStyle := lipgloss.NewStyle().
		Height(contentHeight).
		Width(shared.GlobalState.ScreenWidth)

	content := contentStyle.Render(contentView)

	return lipgloss.JoinVertical(lipgloss.Left, content, tabBarView)
}

func (manager *manager) SwitchScreen(screen types.Screen) types.Command {
	return nil
}

func (manager *manager) OnKeyPress(key types.KeyPress) types.Command {
	switch key.Type {
	case types.CtrlC:
		return manager.DispatchEvent(types.Quit{})
	case types.KeyTab:
		manager.nextTab()
	case types.KeyShiftTab:
		manager.prevTab()
	case types.Alt0:
		manager.switchToTab(0)
	case types.Alt1:
		manager.switchToTab(1)
	case types.Alt2:
		manager.switchToTab(2)
	case types.Alt3:
		manager.switchToTab(3)
	case types.Alt4:
		manager.switchToTab(4)
	case types.Alt5:
		manager.switchToTab(5)
	case types.Alt6:
		manager.switchToTab(6)
	case types.Alt7:
		manager.switchToTab(7)
	case types.Alt8:
		manager.switchToTab(8)
	case types.Alt9:
		manager.switchToTab(9)
	}

	return nil
}

func (manager *manager) DispatchEvent(event types.Event) types.Command {
	return func() types.Event {
		return event
	}
}

func (manager *manager) nextTab() {
	count := manager.tabBar.Count()
	if count > 1 {
		newIndex := (manager.tabBar.ActiveIndex() + 1) % count
		manager.tabBar.SetActive(newIndex)
	}
}

func (manager *manager) prevTab() {
	count := manager.tabBar.Count()
	if count > 1 {
		newIndex := manager.tabBar.ActiveIndex() - 1
		if newIndex < 0 {
			newIndex = count - 1
		}
		manager.tabBar.SetActive(newIndex)
	}
}

func (manager *manager) switchToTab(index int) {
	if index < manager.tabBar.Count() {
		manager.tabBar.SetActive(index)
	}
}

func (manager *manager) updateActiveScreen(screen types.Screen) {
	manager.tabBar.UpdateActiveScreen(screen)
}
