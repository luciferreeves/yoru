package screens

import (
	"yoru/screens/components"
	"yoru/shared"
	"yoru/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var ScreenManager = &manager{
	tabBar: components.TabBar,
}

func (manager *manager) Init() tea.Cmd {
	manager.tabBar.AddTab(types.Tab{
		Name:   "Home",
		Screen: homeScreen,
	})

	return homeScreen.Init()
}

func (manager *manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch message := msg.(type) {
	case types.AddTabMsg:
		// Add new tab and switch to it
		manager.tabBar.AddTab(types.Tab{
			Name:   message.TabName,
			Screen: message.Screen,
		})
		manager.tabBar.SwitchToLastTab()
		// Initialize the new screen
		return manager, message.Screen.Init()
	case types.CloseTabMsg:
		manager.tabBar.RemoveCurrentTab()
		return manager, nil
	case tea.KeyMsg:
		// Check if current screen is in terminal key capture mode
		screen := manager.tabBar.GetCurrentScreen()
		if termScreen, ok := screen.(*terminalScreen); ok {
			if termScreen.GetKeyCaptureMode() == types.KeyCaptureTerminal {
				// In terminal mode, pass ALL keys to terminal screen
				// (terminal.go handles Ctrl+] to release capture)
				current, command := screen.Update(msg)
				manager.tabBar.UpdateCurrentScreen(current)
				return manager, command
			}
		}

		// In normal mode, handle global keys
		if command := manager.OnKeyPress(message); command != nil {
			return manager, command
		}
	case tea.WindowSizeMsg:
		shared.GlobalState.ScreenWidth = message.Width
		shared.GlobalState.ScreenHeight = message.Height
		// Pass window size to current screen too
		screen := manager.tabBar.GetCurrentScreen()
		if screen != nil {
			current, command := screen.Update(msg)
			manager.tabBar.UpdateCurrentScreen(current)
			return manager, command
		}
		return manager, nil
	}

	screen := manager.tabBar.GetCurrentScreen()
	if screen != nil {
		current, command := screen.Update(msg)
		manager.tabBar.UpdateCurrentScreen(current)
		return manager, command
	}

	return manager, nil
}

func (manager *manager) View() string {
	activeScreen := manager.tabBar.GetCurrentScreen()

	var contentView string
	if activeScreen != nil {
		contentView = activeScreen.View()
	}

	tabBarView := manager.tabBar.Render()
	return lipgloss.JoinVertical(lipgloss.Left, contentView, tabBarView)
}

func (manager *manager) SwitchScreen(screen types.Screen) tea.Cmd {
	return nil
}

func (manager *manager) OnKeyPress(key tea.KeyMsg) tea.Cmd {
	switch key.Type {
	case tea.KeyCtrlC, tea.KeyCtrlQ:
		return tea.Quit
	case tea.KeyTab:
		manager.tabBar.NextTab()
	case tea.KeyShiftTab:
		manager.tabBar.PrevTab()
	default:
		if key.Alt {
			switch key.String() {
			case "alt+0":
				manager.tabBar.SwitchToTab(0)
			case "alt+1":
				manager.tabBar.SwitchToTab(1)
			case "alt+2":
				manager.tabBar.SwitchToTab(2)
			case "alt+3":
				manager.tabBar.SwitchToTab(3)
			case "alt+4":
				manager.tabBar.SwitchToTab(4)
			case "alt+5":
				manager.tabBar.SwitchToTab(5)
			case "alt+6":
				manager.tabBar.SwitchToTab(6)
			case "alt+7":
				manager.tabBar.SwitchToTab(7)
			case "alt+8":
				manager.tabBar.SwitchToTab(8)
			case "alt+9":
				manager.tabBar.SwitchToTab(9)
			}
		}
	}

	return nil
}

func (manager *manager) DispatchEvent(event tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return event
	}
}
