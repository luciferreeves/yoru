package types

import tea "github.com/charmbracelet/bubbletea"

type KeyType string
type KeyPress struct {
	Key  tea.KeyMsg
	Type KeyType
}

const (
	// Alt + digits
	Alt0 KeyType = "alt+0"
	Alt1 KeyType = "alt+1"
	Alt2 KeyType = "alt+2"
	Alt3 KeyType = "alt+3"
	Alt4 KeyType = "alt+4"
	Alt5 KeyType = "alt+5"
	Alt6 KeyType = "alt+6"
	Alt7 KeyType = "alt+7"
	Alt8 KeyType = "alt+8"
	Alt9 KeyType = "alt+9"

	// Ctrl
	CtrlC KeyType = "ctrl+c"
	CtrlN KeyType = "ctrl+n"

	// Special
	KeyEnter    KeyType = "enter"
	KeyEscape   KeyType = "escape"
	KeyUp       KeyType = "up"
	KeyDown     KeyType = "down"
	KeyLeft     KeyType = "left"
	KeyRight    KeyType = "right"
	KeyTab      KeyType = "tab"
	KeyShiftTab KeyType = "shift+tab"
)
