package screens

import (
	"yoru/shared"
	"yoru/types"
)

var ScreenManager types.ScreenManager = &screenManager{
	Current: _root(),
}

func (manager *screenManager) Init() types.Command {
	return manager.Current.Init()
}

func (manager *screenManager) Update(event types.Event) (types.Screen, types.Command) {
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

func (manager *screenManager) View() string {
	return manager.Current.View()
}

func (manager *screenManager) SwitchScreen(screen types.Screen) types.Command {
	return manager.DispatchEvent(types.ScreenSwitched{Screen: screen})
}

func (manager *screenManager) OnKeyPress(key types.KeyPress) types.Command {
	switch key.Type {
	case types.CtrlC:
		return manager.DispatchEvent(types.Quit{})
	default:
		return nil
	}
}

func (manager *screenManager) DispatchEvent(event types.Event) types.Command {
	return func() types.Event {
		return event
	}
}
