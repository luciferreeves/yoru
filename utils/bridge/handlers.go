package bridge

import (
	"strings"
	"unicode"
	"yoru/types"

	tea "github.com/charmbracelet/bubbletea"
)

func (adapter *bridge) convertMessage(msg tea.Msg) types.Event {
	switch teaMsg := msg.(type) {
	case tea.KeyMsg:
		return types.KeyPress{
			Key:  teaMsg,
			Type: adapter.normalizeKey(teaMsg),
		}
	default:
		return msg
	}
}

func (adapter *bridge) normalizeKey(msg tea.KeyMsg) types.KeyType {
	if msg.Type >= tea.KeyCtrlA && msg.Type <= tea.KeyCtrlZ {
		letter := rune('a' + (msg.Type - tea.KeyCtrlA))
		return types.KeyType("ctrl+" + string(letter))
	}

	switch msg.Type {
	case tea.KeyEnter:
		return "enter"
	case tea.KeyEscape:
		return "escape"
	case tea.KeyUp:
		return "up"
	case tea.KeyDown:
		return "down"
	case tea.KeyLeft:
		return "left"
	case tea.KeyRight:
		return "right"
	case tea.KeyTab:
		return "tab"
	case tea.KeyShiftTab:
		return "shift+tab"
	}

	if len(msg.Runes) == 0 {
		return ""
	}

	r := msg.Runes[0]
	base := string(unicode.ToLower(r))
	parts := []string{}

	if msg.Alt {
		parts = append(parts, "alt")
	}

	if unicode.IsUpper(r) {
		parts = append(parts, "shift")
	}

	parts = append(parts, base)

	return types.KeyType(strings.Join(parts, "+"))
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
