package types

import tea "github.com/charmbracelet/bubbletea"

type Screen interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Screen, tea.Cmd)
	View() string
}

type ScreenManager interface {
	tea.Model
	SwitchScreen(screen Screen) tea.Cmd
	OnKeyPress(key tea.KeyMsg) tea.Cmd
	DispatchEvent(event tea.Msg) tea.Cmd
}
