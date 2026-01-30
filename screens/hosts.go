package screens

import (
	"yoru/models"
	"yoru/repository"
	"yoru/screens/components"
	"yoru/screens/forms"
	"yoru/screens/popups"
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
	deletePopup: popups.NewDeleteHostPopup(),
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
	if screen.deletePopup.IsVisible() {
		screen.deletePopup.Update(msg)
		return screen, nil
	}

	switch message := msg.(type) {
	case tea.KeyMsg:
		switch message.Type {
		case tea.KeyEnter:
			if screen.focusedArea == sidebarFocus {
				if cmd := screen.OnKeyPress(message); cmd != nil {
					return screen, cmd
				}
				return screen, nil
			}
		case tea.KeyEscape:
			if screen.focusedArea == formFocus {
				if cmd := screen.OnKeyPress(message); cmd != nil {
					return screen, cmd
				}
				return screen, nil
			}
		case tea.KeyCtrlN:
			if cmd := screen.OnKeyPress(message); cmd != nil {
				return screen, cmd
			}
			return screen, nil
		}

		if screen.focusedArea == sidebarFocus && (message.String() == "d" || message.String() == "D") {
			selectedHost := screen.sidebar.GetSelected()
			if selectedHost != nil {
				screen.deletePopup.Show(
					selectedHost.Name,
					func(dontAskAgain bool) {
						if err := repository.DeleteHost(selectedHost.ID); err == nil {
							allHosts, _ := repository.GetAllHosts()
							screen.sidebar.SetHosts(allHosts)
							if len(allHosts) > 0 {
								selected := screen.sidebar.GetSelected()
								if selected != nil {
									screen.form.LoadHost(selected)
								}
							}
						}
					},
					func() {},
				)
			}
			return screen, nil
		}
	}

	// Block left/right navigation when popup is visible
	if screen.deletePopup.IsVisible() {
		return screen, nil
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

	content := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, form)

	if screen.deletePopup.IsVisible() {
		popupView := screen.deletePopup.Render()
		return popupView
	}

	return content
}

func (screen *hosts) OnKeyPress(key tea.KeyMsg) tea.Cmd {
	switch key.Type {
	case tea.KeyEnter:
		if screen.focusedArea == sidebarFocus {
			screen.filterWasActive = screen.sidebar.IsFilterActive()
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
			screen.sidebar.SetFilterActive(screen.filterWasActive)
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
