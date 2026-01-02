package screens

import "yoru/types"

var homeScreen = &home{}

func (screen *home) Init() types.Command {
	return nil
}

func (screen *home) Update(event types.Event) (types.Screen, types.Command) {
	return screen, nil
}

func (screen *home) View() string {
	return "Home Screen"
}
