package screens

import (
	"yoru/shared"
	"yoru/types"
)

var ScreenManager = &manager{
	Current: rootScreen,
}

func (manager *manager) Init() types.Command {
	return manager.Current.Init()
}

func (manager *manager) Update(event types.Event) (types.Screen, types.Command) {
	switch message := event.(type) {
	case types.KeyPress:
		if command := manager.OnKeyPress(message); command != nil {
			return manager, command
		}
	case types.WindowResized:
		shared.GlobalState.ScreenWidth, shared.GlobalState.ScreenHeight = message.Width, message.Height
		return manager, nil
	case types.ScreenSwitched:
		manager.Current = message.Screen
		return manager, manager.Current.Init()
	}

	current, command := manager.Current.Update(event)
	manager.Current = current
	return manager, command
}

func (manager *manager) View() string {
	return manager.Current.View()
}

func (manager *manager) SwitchScreen(screen types.Screen) types.Command {
	return manager.DispatchEvent(types.ScreenSwitched{Screen: screen})
}

func (manager *manager) OnKeyPress(key types.KeyPress) types.Command {
	switch key.Type {
	case types.CtrlC:
		return manager.DispatchEvent(types.Quit{})
	default:
		return nil
	}
}

func (manager *manager) DispatchEvent(event types.Event) types.Command {
	return func() types.Event {
		return event
	}
}
