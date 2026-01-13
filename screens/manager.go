package screens

import (
	"yoru/screens/components"
	"yoru/shared"
	"yoru/types"

	"github.com/charmbracelet/lipgloss"
)

var ScreenManager = &manager{
	tabBar: components.TabBar,
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

	screen := manager.tabBar.GetCurrentScreen()
	if screen != nil {
		current, command := screen.Update(event)
		manager.tabBar.UpdateCurrentScreen(current)
		return manager, command
	}

	return manager, nil
}

func (manager *manager) View() string {
	activeScreen := manager.tabBar.GetCurrentScreen()
	tabBarView := manager.tabBar.Render()

	var contentView string
	if activeScreen != nil {
		contentView = activeScreen.View()
	}

	return lipgloss.JoinVertical(lipgloss.Left, contentView, tabBarView)
}

func (manager *manager) SwitchScreen(screen types.Screen) types.Command {
	return nil
}

func (manager *manager) OnKeyPress(key types.KeyPress) types.Command {
	switch key.Type {
	case types.CtrlC:
		return manager.DispatchEvent(types.Quit{})
	case types.KeyTab:
		manager.tabBar.NextTab()
	case types.KeyShiftTab:
		manager.tabBar.PrevTab()
	case types.Alt0:
		manager.tabBar.SwitchToTab(0)
	case types.Alt1:
		manager.tabBar.SwitchToTab(1)
	case types.Alt2:
		manager.tabBar.SwitchToTab(2)
	case types.Alt3:
		manager.tabBar.SwitchToTab(3)
	case types.Alt4:
		manager.tabBar.SwitchToTab(4)
	case types.Alt5:
		manager.tabBar.SwitchToTab(5)
	case types.Alt6:
		manager.tabBar.SwitchToTab(6)
	case types.Alt7:
		manager.tabBar.SwitchToTab(7)
	case types.Alt8:
		manager.tabBar.SwitchToTab(8)
	case types.Alt9:
		manager.tabBar.SwitchToTab(9)
	}

	return nil
}

func (manager *manager) DispatchEvent(event types.Event) types.Command {
	return func() types.Event {
		return event
	}
}
