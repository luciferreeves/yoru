package types

import tea "github.com/charmbracelet/bubbletea"

type KeyPress struct {
	Key  tea.KeyMsg
	Type rune
}

const (
	CtrlC     rune = 0x03
	KeyW      rune = 'w'
	KeyShiftW rune = 'W'
)
