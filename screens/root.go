package screens

import "yoru/types"

var rootScreen = &root{
	toggled: false,
}

func (screen *root) Init() types.Command {
	return nil
}

func (screen *root) Update(event types.Event) (types.Screen, types.Command) {
	switch message := event.(type) {
	case types.KeyPress:
		return screen, screen.onKeyPress(message)
	}

	return screen, nil
}

func (screen *root) View() string {
	if screen.toggled {
		return "Well done!"
	}

	return "Welcome"
}

func (screen *root) onKeyPress(key types.KeyPress) types.Command {
	switch key.Type {
	case types.KeyW, types.KeyShiftW:
		screen.toggled = !screen.toggled
		return nil
	default:
		return nil
	}
}
