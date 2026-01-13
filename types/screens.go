package types

type ProgramOption int

const (
	WithAltScreen ProgramOption = iota
)

type Command func() Event

type Screen interface {
	Init() Command
	Update(msg Event) (Screen, Command)
	View() string
}

type ScreenManager interface {
	Screen
	SwitchScreen(screen Screen) Command
	OnKeyPress(key KeyPress) Command
	DispatchEvent(event Event) Command
}
