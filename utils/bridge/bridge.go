package bridge

import (
	"fmt"
	"os"
	"yoru/shared"
	"yoru/types"

	tea "github.com/charmbracelet/bubbletea"
)

func New(manager types.ScreenManager, options ...types.ProgramOption) error {
	adapter := &bridge{manager: manager}
	program := tea.NewProgram(adapter, adapter.convertOptions(options...)...)
	_, err := program.Run()
	return err
}

func ExitOnError(err error) {
	fmt.Printf("An error occurred. %s will now exit. Error: %v\n", shared.PrettyName, err)
	os.Exit(1)
}
