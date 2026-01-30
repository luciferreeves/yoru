package components

import (
	"yoru/screens/styles"
	"yoru/shared"
	"yoru/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Popup struct {
	visible   bool
	width     int
	maxHeight int
	content   string
	onUpdate  func(tea.Msg) bool
}

func NewPopup() *Popup {
	return &Popup{
		visible:   false,
		width:     60,
		maxHeight: 20,
	}
}

func (popup *Popup) Show(content string, onUpdate func(tea.Msg) bool) {
	popup.visible = true
	popup.content = content
	popup.onUpdate = onUpdate
}

func (popup *Popup) Hide() {
	popup.visible = false
}

func (popup *Popup) IsVisible() bool {
	return popup.visible
}

func (popup *Popup) Update(msg tea.Msg) bool {
	if !popup.visible {
		return false
	}

	if popup.onUpdate != nil {
		return popup.onUpdate(msg)
	}
	return false
}

func (popup *Popup) SetContent(content string) {
	popup.content = content
}

func (popup *Popup) Render() string {
	if !popup.visible {
		return ""
	}

	popupBox := styles.PopupContainer.
		Width(popup.width).
		MaxHeight(popup.maxHeight).
		Render(popup.content)

	availHeight := shared.GlobalState.ScreenHeight - 4
	if availHeight < 1 {
		availHeight = 1
	}

	return lipgloss.Place(
		shared.GlobalState.ScreenWidth,
		availHeight,
		lipgloss.Center,
		lipgloss.Center,
		popupBox,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color(types.Crust)),
	)
}
