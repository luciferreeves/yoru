package popups

import (
	"yoru/screens/components"
	"yoru/screens/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DeleteKeychainPopup struct {
	popup          *components.Popup
	itemName       string
	itemType       string
	checkboxValue  bool
	selectedButton int // 0 = No, 1 = Yes
	onConfirm      func(dontAskAgain bool)
	onCancel       func()
}

func NewDeleteKeychainPopup() *DeleteKeychainPopup {
	dkp := &DeleteKeychainPopup{
		popup:         components.NewPopup(),
		checkboxValue: false,
	}
	return dkp
}

func (dkp *DeleteKeychainPopup) Show(itemName string, itemType string, onConfirm func(bool), onCancel func()) {
	dkp.itemName = itemName
	dkp.itemType = itemType
	dkp.onConfirm = onConfirm
	dkp.onCancel = onCancel
	dkp.checkboxValue = false
	dkp.selectedButton = 0 // Default to No

	dkp.popup.Show(dkp.buildContent(), dkp.handleInput)
}

func (dkp *DeleteKeychainPopup) Hide() {
	dkp.popup.Hide()
}

func (dkp *DeleteKeychainPopup) IsVisible() bool {
	return dkp.popup.IsVisible()
}

func (dkp *DeleteKeychainPopup) Update(msg tea.Msg) {
	dkp.popup.Update(msg)
}

func (dkp *DeleteKeychainPopup) Render() string {
	return dkp.popup.Render()
}

func (dkp *DeleteKeychainPopup) handleInput(msg tea.Msg) bool {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case " ":
			dkp.checkboxValue = !dkp.checkboxValue
			dkp.popup.SetContent(dkp.buildContent())
			return true
		case "left":
			dkp.selectedButton = 1 // Yes
			dkp.popup.SetContent(dkp.buildContent())
			return true
		case "right":
			dkp.selectedButton = 0 // No
			dkp.popup.SetContent(dkp.buildContent())
			return true
		case "enter":
			if dkp.selectedButton == 1 {
				if dkp.onConfirm != nil {
					dkp.onConfirm(dkp.checkboxValue)
				}
			} else {
				if dkp.onCancel != nil {
					dkp.onCancel()
				}
			}
			dkp.Hide()
			return true
		case "y", "Y":
			if dkp.onConfirm != nil {
				dkp.onConfirm(dkp.checkboxValue)
			}
			dkp.Hide()
			return true
		case "n", "N", "esc":
			if dkp.onCancel != nil {
				dkp.onCancel()
			}
			dkp.Hide()
			return true
		}
	}
	return false
}

func (dkp *DeleteKeychainPopup) buildContent() string {
	title := styles.PopupTitle.Render("Delete " + dkp.itemType)
	message := styles.PopupMessage.Render("Are you sure you want to delete \"" + dkp.itemName + "\"?")

	checkboxIcon := "[ ]"
	checkboxStyle := styles.PopupCheckbox
	if dkp.checkboxValue {
		checkboxIcon = "[x]"
		checkboxStyle = styles.PopupCheckboxChecked
	}
	checkbox := checkboxStyle.Render(checkboxIcon + " Never ask this again")

	yesPrefix := "  "
	noPrefix := "  "
	if dkp.selectedButton == 1 {
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
