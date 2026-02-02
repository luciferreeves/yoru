package styles

import (
	"yoru/types"

	"github.com/charmbracelet/lipgloss"
)

var (
	PopupOverlay = lipgloss.NewStyle().
			Background(lipgloss.Color(types.Crust)).
			Foreground(lipgloss.Color(types.Text))

	PopupContainer = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(types.Lavender)).
			Background(lipgloss.Color(types.Base)).
			Padding(1, 2).
			MaxHeight(20)

	PopupTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Lavender)).
			Bold(true).
			MarginBottom(1)

	PopupMessage = lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Text)).
			MarginBottom(1)

	PopupCheckbox = lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Subtext0)).
			MarginBottom(1)

	PopupCheckboxChecked = lipgloss.NewStyle().
				Foreground(lipgloss.Color(types.Lavender)).
				MarginBottom(1)

	PopupButtonNo = lipgloss.NewStyle().
			Background(lipgloss.Color(types.Red)).
			Foreground(lipgloss.Color(types.Base)).
			Padding(0, 2).
			Bold(true).
			Align(lipgloss.Center)

	PopupButtonYes = lipgloss.NewStyle().
			Background(lipgloss.Color(types.Surface0)).
			Foreground(lipgloss.Color(types.Text)).
			Padding(0, 2).
			Bold(true).
			Align(lipgloss.Center)

	PopupButtonSelected = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color(types.Lavender)).
				BorderTop(false).
				BorderLeft(false).
				BorderRight(false).
				BorderBottom(true)

	PopupButtonsContainer = lipgloss.NewStyle().
				MarginTop(1)

	PopupItemSelected = lipgloss.NewStyle().
				Foreground(lipgloss.Color(types.Base)).
				Background(lipgloss.Color(types.Lavender)).
				Bold(true)

	PopupItemNormal = lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Text)).
			Background(lipgloss.Color(types.Base))
)
