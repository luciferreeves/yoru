package screens

import "yoru/types"

func _root() types.Screen {
	return &rootScreen{}
}

func (root *rootScreen) Init() types.Command {
	return nil
}

func (root *rootScreen) Update(event types.Event) (types.Screen, types.Command) {
	switch message := event.(type) {
	case types.KeyPress:
		return root, root.onKeyPress(message)
	}

	return root, nil
}

func (root *rootScreen) View() string {
	if root.toggled {
		return "Well done!"
	}

	return "Welcome"
}

func (root *rootScreen) onKeyPress(key types.KeyPress) types.Command {
	switch key.Type {
	case types.KeyW, types.KeyShiftW:
		root.toggled = !root.toggled
		return nil
	default:
		return nil
	}
}
