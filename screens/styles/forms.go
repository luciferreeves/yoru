package styles

import (
	"yoru/types"

	"github.com/charmbracelet/lipgloss"
)

var (
	FormLabel = lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Mauve)).
			Bold(true).
			Width(12)

	FormText = lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Text))

	FormTextFocused = lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Mauve))

	FormPlaceholder = lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Subtext0))

	FormEmpty = lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Subtext0)).
			Padding(2)

	FormError = lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Red))
)
