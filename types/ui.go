package types

type Tab struct {
	Name   string
	Screen Screen
}

type TabBar interface {
	AddTab(tab Tab)
	GetCurrentScreen() Screen
	UpdateCurrentScreen(screen Screen)
	SwitchToTab(index int)
	NextTab()
	PrevTab()
	Render() string
}

type NavBar interface {
	GetActiveTab() (int, string)
	SwitchToTab(index int)
	NextTab()
	PrevTab()
	Render() string
}
