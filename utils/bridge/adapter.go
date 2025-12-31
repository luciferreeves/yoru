package bridge

import (
	"yoru/types"

	tea "github.com/charmbracelet/bubbletea"
)

func (adapter *bridge) Init() tea.Cmd {
	cmd := adapter.manager.Init()
	if cmd == nil {
		return nil
	}
	return adapter.wrapCommand(cmd)
}

func (adapter *bridge) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	event := adapter.convertMessage(msg)
	screen, cmd := adapter.manager.Update(event)
	adapter.manager = screen.(types.ScreenManager)
	if cmd == nil {
		return adapter, nil
	}
	return adapter, adapter.handleCommand(cmd)
}

func (adapter *bridge) View() string {
	return adapter.manager.View()
}
