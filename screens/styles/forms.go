package styles

import (
	"yoru/types"

	"github.com/charmbracelet/lipgloss"
)

var (
	FormContainer = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(types.Lavender)).
			Padding(2, 4)

	FormTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Lavender)).
			Bold(true).
			MarginBottom(1)

	FormSectionTitle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(types.Peach)).
				Bold(true).
				MarginTop(1).
				MarginBottom(1)

	FormLabel = lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Mauve)).
			Bold(true).
			Width(12)

	FormFieldContainer = lipgloss.NewStyle().
				MarginBottom(1)

	FormInput = lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Text))

	FormInputFocused = lipgloss.NewStyle().
				Foreground(lipgloss.Color(types.Lavender))

	FormCheckbox = lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Text)).
			MarginRight(1)

	FormCheckboxFocused = lipgloss.NewStyle().
				Foreground(lipgloss.Color(types.Lavender)).
				Bold(true).
				MarginRight(1)

	FormCheckboxLabel = lipgloss.NewStyle().
				Foreground(lipgloss.Color(types.Text))

	FormCheckboxLabelFocused = lipgloss.NewStyle().
					Foreground(lipgloss.Color(types.Lavender))

	FormText = lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Text))

	FormTextFocused = lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Mauve))

	FormPlaceholder = lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Subtext0)).
			Italic(true)

	FormEmpty = lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Subtext0)).
			Italic(true)

	FormError = lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Red)).
			MarginLeft(12)
)
