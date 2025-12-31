package bridge

import (
	"yoru/types"

	tea "github.com/charmbracelet/bubbletea"
)

func (adapter *bridge) convertMessage(msg tea.Msg) types.Event {
	switch teaMsg := msg.(type) {
	case tea.KeyMsg:
		return types.KeyPress{
			Key:  teaMsg,
			Type: append(teaMsg.Runes, rune(teaMsg.Type))[0],
		}
	default:
		return msg
	}
}

func (adapter *bridge) wrapCommand(cmd types.Command) tea.Cmd {
	return func() tea.Msg { return cmd() }
}

func (adapter *bridge) handleCommand(cmd types.Command) tea.Cmd {
	result := cmd()
	if _, ok := result.(types.Quit); ok {
		return tea.Quit
	}
	return adapter.wrapCommand(func() types.Event { return result })
}

func (adapter *bridge) convertOptions(options ...types.ProgramOption) []tea.ProgramOption {
	converted := make([]tea.ProgramOption, 0, len(options))
	for _, option := range options {
		switch option {
		case types.WithAltScreen:
			converted = append(converted, tea.WithAltScreen())
		}
	}
	return converted
}
