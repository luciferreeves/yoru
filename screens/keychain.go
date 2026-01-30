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

type keychain struct {
	sidebar     *components.KeychainSidebar
	form        *forms.KeychainForm
	focusedArea focusArea
	deletePopup *popups.DeleteKeychainPopup
}

const (
	keychainSidebarFocus focusArea = iota
	keychainFormFocus
)

const (
	keychainSidebarWidth = 30
	keychainFormHeight   = 4
)

var keychainScreen = &keychain{
	sidebar:     components.NewKeychainSidebar(),
	form:        forms.NewKeychainForm(),
	focusedArea: keychainSidebarFocus,
	deletePopup: popups.NewDeleteKeychainPopup(),
}

func (screen *keychain) Init() tea.Cmd {
	keys, _ := repository.GetAllKeys()
	identities, _ := repository.GetAllIdentities()
	screen.sidebar.SetItems(keys, identities)

	if len(keys) > 0 {
		screen.form.LoadKey(&keys[0])
	} else if len(identities) > 0 {
		screen.form.LoadIdentity(&identities[0])
	}

	return nil
}

func (screen *keychain) Update(msg tea.Msg) (types.Screen, tea.Cmd) {
	if screen.deletePopup.IsVisible() {
		screen.deletePopup.Update(msg)
		return screen, nil
	}

	switch message := msg.(type) {
	case tea.KeyMsg:
		switch message.Type {
		case tea.KeyEnter:
			if screen.focusedArea == keychainSidebarFocus {
				if cmd := screen.OnKeyPress(message); cmd != nil {
					return screen, cmd
				}
				return screen, nil
			}
		case tea.KeyEscape:
			if screen.focusedArea == keychainFormFocus {
				if screen.form.IsTextAreaEditing() {
					screen.form.Update(msg)
					return screen, nil
				}
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

		if screen.focusedArea == keychainSidebarFocus && (message.String() == "d" || message.String() == "D") {
			selectedItem := screen.sidebar.GetSelected()
			if selectedItem != nil {
				screen.deletePopup.Show(
					selectedItem.Name,
					selectedItem.ItemType,
					func(dontAskAgain bool) {
						var err error
						if selectedItem.ItemType == "Key" {
							err = repository.DeleteKey(selectedItem.ID)
						} else {
							err = repository.DeleteIdentity(selectedItem.ID)
						}

						if err == nil {
							keys, _ := repository.GetAllKeys()
							identities, _ := repository.GetAllIdentities()
							screen.sidebar.SetItems(keys, identities)

							newSelected := screen.sidebar.GetSelected()
							if newSelected != nil {
								if newSelected.ItemType == "Key" {
									for _, key := range keys {
										if key.ID == newSelected.ID {
											screen.form.LoadKey(&key)
											break
										}
									}
								} else {
									for _, identity := range identities {
										if identity.ID == newSelected.ID {
											screen.form.LoadIdentity(&identity)
											break
										}
									}
								}
							} else {
								screen.form.Clear()
							}
						}
					},
					func() {},
				)
			}
			return screen, nil
		}
	}

	if screen.deletePopup.IsVisible() {
		return screen, nil
	}

	switch screen.focusedArea {
	case keychainSidebarFocus:
		screen.sidebar.Update(msg)

		selectedItem := screen.sidebar.GetSelected()
		if selectedItem != nil {
			if selectedItem.ItemType != screen.form.GetLastSelectedType() || selectedItem.ID != screen.form.GetLastSelectedID() {
				screen.form.Save()

				keys, _ := repository.GetAllKeys()
				identities, _ := repository.GetAllIdentities()

				if selectedItem.ItemType == "Key" {
					for _, key := range keys {
						if key.ID == selectedItem.ID {
							screen.form.LoadKey(&key)
							break
						}
					}
				} else {
					for _, identity := range identities {
						if identity.ID == selectedItem.ID {
							screen.form.LoadIdentity(&identity)
							break
						}
					}
				}
			}
		}
	case keychainFormFocus:
		screen.form.Update(msg)
	}

	return screen, nil
}

func (screen *keychain) View() string {
	sidebarView := screen.sidebar.Render()
	formView := screen.form.Render()

	sidebar := styles.SidebarArea.
		Width(keychainSidebarWidth).
		Render(sidebarView)

	formWidth := shared.GlobalState.ScreenWidth - keychainSidebarWidth - 6
	formAreaHeight := shared.GlobalState.ScreenHeight - 8

	centeredForm := lipgloss.Place(
		formWidth,
		formAreaHeight,
		lipgloss.Center,
		lipgloss.Center,
		formView,
	)

	form := styles.FormArea.
		Width(formWidth).
		Height(formAreaHeight).
		Render(centeredForm)

	content := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, form)

	if screen.deletePopup.IsVisible() {
		popupView := screen.deletePopup.Render()
		return popupView
	}
	return content
}

func (screen *keychain) OnKeyPress(key tea.KeyMsg) tea.Cmd {
	if key.Type == tea.KeyCtrlN {
		newIdentity := &models.Identity{
			Name:     "New Identity",
			Username: "",
			Password: "",
		}
		if err := repository.CreateIdentity(newIdentity); err == nil {
			keys, _ := repository.GetAllKeys()
			identities, _ := repository.GetAllIdentities()
			screen.sidebar.SetItems(keys, identities)

			screen.form.LoadIdentity(newIdentity)

			screen.sidebar.SelectItemByID(newIdentity.ID, "Identity")

			screen.focusedArea = keychainFormFocus
			screen.form.SetFocused(true)
		}
		return nil
	}

	if key.Type == tea.KeyEnter && screen.focusedArea == keychainSidebarFocus {
		if !screen.sidebar.IsFilterActive() {
			screen.focusedArea = keychainFormFocus
			screen.form.SetFocused(true)
		}
		return nil
	}

	if key.Type == tea.KeyEscape && screen.focusedArea == keychainFormFocus {
		screen.form.Save()
		screen.form.SetFocused(false)
		screen.focusedArea = keychainSidebarFocus

		keys, _ := repository.GetAllKeys()
		identities, _ := repository.GetAllIdentities()
		screen.sidebar.SetItems(keys, identities)
		return nil
	}

	return nil
}
