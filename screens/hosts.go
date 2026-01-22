package screens

import (
	"yoru/models"
	"yoru/repository"
	"yoru/screens/components"
	"yoru/screens/forms"
	"yoru/screens/styles"
	"yoru/shared"
	"yoru/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	sidebarFocus focusArea = iota
	formFocus
)

const (
	sidebarWidth = 30
	formHeight   = 4
)

var hostsScreen = &hosts{
	sidebar:     components.NewHostsSidebar(),
	form:        forms.NewHostForm(),
	focusedArea: sidebarFocus,
}

func (screen *hosts) Init() tea.Cmd {
	allHosts, _ := repository.GetAllHosts()
	screen.sidebar.SetHosts(allHosts)

	if len(allHosts) > 0 {
		screen.form.LoadHost(&allHosts[0])
	}

	return nil
}

func (screen *hosts) Update(msg tea.Msg) (types.Screen, tea.Cmd) {
	switch message := msg.(type) {
	case tea.KeyMsg:
		switch message.Type {
		case tea.KeyEnter, tea.KeyEscape, tea.KeyCtrlN:
			if cmd := screen.OnKeyPress(message); cmd != nil {
				return screen, cmd
			}
		}
	}

	switch screen.focusedArea {
	case sidebarFocus:
		screen.sidebar.Update(msg)

		selectedHost := screen.sidebar.GetSelected()
		if selectedHost != nil && screen.form.GetLastSelectedHostID() != selectedHost.ID {
			screen.form.Save()
			screen.form.LoadHost(selectedHost)
		}
	case formFocus:
		screen.form.Update(msg)
	}

	return screen, nil
}

func (screen *hosts) View() string {
	sidebarView := screen.sidebar.Render()
	formView := screen.form.Render()

	sidebar := styles.SidebarArea.
		Width(sidebarWidth).
		Render(sidebarView)

	form := styles.FormArea.
		Width(shared.GlobalState.ScreenWidth - sidebarWidth - 6).
		Height(formHeight).
		Render(formView)

	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, form)
}

func (screen *hosts) OnKeyPress(key tea.KeyMsg) tea.Cmd {
	switch key.Type {
	case tea.KeyEnter:
		if screen.focusedArea == sidebarFocus {
			screen.focusedArea = formFocus
			screen.form.SetFocused(true)
			return nil
		}
		return nil

	case tea.KeyEscape:
		if screen.focusedArea == formFocus {
			screen.focusedArea = sidebarFocus
			screen.form.SetFocused(false)
			screen.form.Save()
			return nil
		}
		return nil

	case tea.KeyCtrlN:
		newHost := &models.Host{
			Name:           "New Host",
			Hostname:       "0.0.0.0",
			Mode:           types.ModeSSH,
			Port:           22,
			CredentialID:   1,
			CredentialType: types.CredentialIdentity,
		}
		if err := repository.CreateHost(newHost); err == nil {
			allHosts, _ := repository.GetAllHosts()
			screen.sidebar.SetHosts(allHosts)
			if len(allHosts) > 0 {
				screen.form.LoadHost(&allHosts[0])
			}
			screen.focusedArea = formFocus
			screen.form.SetFocused(true)
		}
		return nil
	}

	return nil
}
