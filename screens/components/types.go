package components

import "yoru/types"

type tabBar struct {
	types.TabBar
	tabs        []types.Tab
	activeIndex int
}

type navBar struct {
	types.NavBar
	items       []string
	activeIndex int
}
