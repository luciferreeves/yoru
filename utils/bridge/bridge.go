package bridge

import (
	"yoru/types"

	tea "github.com/charmbracelet/bubbletea"
)

func New(manager types.ScreenManager, options ...types.ProgramOption) error {
	adapter := &bridge{manager: manager}
	program := tea.NewProgram(adapter, adapter.convertOptions(options...)...)
	_, err := program.Run()
	return err
}
