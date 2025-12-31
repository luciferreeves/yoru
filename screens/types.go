package screens

import "yoru/types"

type screenManager struct {
	Current types.Screen
}

type rootScreen struct {
	toggled bool
}
