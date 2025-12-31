package types

import tea "github.com/charmbracelet/bubbletea"

type Event interface {
	tea.Msg
}

type WindowResized tea.WindowSizeMsg

type ScreenSwitched struct {
	Screen Screen
}

type Quit struct{}
