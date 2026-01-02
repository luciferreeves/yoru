package styles

import (
	"yoru/types"

	"github.com/charmbracelet/lipgloss"
)

var (
	ActiveTab = lipgloss.NewStyle().
			Background(lipgloss.Color(types.Mauve)).
			Foreground(lipgloss.Color(types.Base)).
			Padding(0, 3).
			Bold(true)

	InactiveTab = lipgloss.NewStyle().
			Background(lipgloss.Color(types.Surface1)).
			Foreground(lipgloss.Color(types.Subtext0)).
			Padding(0, 3)

	TabBarBackground = lipgloss.NewStyle().
				Background(lipgloss.Color(types.Surface0))

	ContentBackground = lipgloss.NewStyle().
				Background(lipgloss.Color(types.Base)).
				Foreground(lipgloss.Color(types.Text))
)
