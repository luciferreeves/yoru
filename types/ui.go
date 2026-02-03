package types

type Tab struct {
	Name   string
	Screen Screen
}

// AddTabMsg is a message to add a new tab and switch to it
type AddTabMsg struct {
	TabName string
	Screen  Screen
}

// CloseTabMsg is a message to remove the current tab
type CloseTabMsg struct{}

type TabBar interface {
	AddTab(tab Tab)
	RemoveCurrentTab()
	GetCurrentScreen() Screen
	UpdateCurrentScreen(screen Screen)
	SwitchToTab(index int)
	SwitchToLastTab()
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
