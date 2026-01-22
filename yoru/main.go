package main

import (
	"yoru/screens"
	"yoru/utils/errors"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	program := tea.NewProgram(screens.ScreenManager, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		errors.ExitOnBridgeFailedStart(err)
	}
}
