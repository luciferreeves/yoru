package popups

import (
	"yoru/screens/components"
	"yoru/screens/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DeleteHostPopup struct {
	popup          *components.Popup
	hostName       string
	checkboxValue  bool
	selectedButton int // 0 = No, 1 = Yes
	onConfirm      func(dontAskAgain bool)
	onCancel       func()
}

func NewDeleteHostPopup() *DeleteHostPopup {
	dhp := &DeleteHostPopup{
		popup:         components.NewPopup(),
		checkboxValue: false,
	}
	return dhp
}

func (dhp *DeleteHostPopup) Show(hostName string, onConfirm func(bool), onCancel func()) {
	dhp.hostName = hostName
	dhp.onConfirm = onConfirm
	dhp.onCancel = onCancel
	dhp.checkboxValue = false
	dhp.selectedButton = 0 // Default to No

	dhp.popup.Show(dhp.buildContent(), dhp.handleInput)
}

func (dhp *DeleteHostPopup) Hide() {
	dhp.popup.Hide()
}

func (dhp *DeleteHostPopup) IsVisible() bool {
	return dhp.popup.IsVisible()
}

func (dhp *DeleteHostPopup) Update(msg tea.Msg) {
	dhp.popup.Update(msg)
}

func (dhp *DeleteHostPopup) Render() string {
	return dhp.popup.Render()
}

func (dhp *DeleteHostPopup) handleInput(msg tea.Msg) bool {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case " ":
			dhp.checkboxValue = !dhp.checkboxValue
			dhp.popup.SetContent(dhp.buildContent())
			return true
		case "left":
			dhp.selectedButton = 1 // Yes
			dhp.popup.SetContent(dhp.buildContent())
			return true
		case "right":
			dhp.selectedButton = 0 // No
			dhp.popup.SetContent(dhp.buildContent())
			return true
		case "enter":
			if dhp.selectedButton == 1 {
				if dhp.onConfirm != nil {
					dhp.onConfirm(dhp.checkboxValue)
				}
			} else {
				if dhp.onCancel != nil {
					dhp.onCancel()
				}
			}
			dhp.Hide()
			return true
		case "y", "Y":
			if dhp.onConfirm != nil {
				dhp.onConfirm(dhp.checkboxValue)
			}
			dhp.Hide()
			return true
		case "n", "N", "esc":
			if dhp.onCancel != nil {
				dhp.onCancel()
			}
			dhp.Hide()
			return true
		}
	}
	return false
}

func (dhp *DeleteHostPopup) buildContent() string {
	title := styles.PopupTitle.Render("Delete Host")
	message := styles.PopupMessage.Render("Are you sure you want to delete \"" + dhp.hostName + "\"?")

	checkboxIcon := "[ ]"
	checkboxStyle := styles.PopupCheckbox
	if dhp.checkboxValue {
		checkboxIcon = "[x]"
		checkboxStyle = styles.PopupCheckboxChecked
	}
	checkbox := checkboxStyle.Render(checkboxIcon + " Never ask this again (Space)")

	// Add selection indicator
	yesPrefix := "  "
	noPrefix := "  "
	if dhp.selectedButton == 1 {
		yesPrefix = "> "
	} else {
		noPrefix = "> "
	}

	yesButton := styles.PopupButtonYes.Render(yesPrefix + "Yes (y)")
	noButton := styles.PopupButtonNo.Render(noPrefix + "No (n)")

	buttons := lipgloss.JoinHorizontal(lipgloss.Top, yesButton, "  ", noButton)
	buttonsContainer := lipgloss.NewStyle().Width(56).Align(lipgloss.Right).Render(buttons)
	buttonsWithMargin := styles.PopupButtonsContainer.Render(buttonsContainer)

	return lipgloss.JoinVertical(lipgloss.Left, title, message, checkbox, buttonsWithMargin)
}
