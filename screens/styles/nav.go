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

	ActiveNavBar = lipgloss.NewStyle().
			Background(lipgloss.Color(types.Lavender)).
			Foreground(lipgloss.Color(types.Base)).
			Padding(0, 2).
			Bold(true)

	InactiveNavBar = lipgloss.NewStyle().
			Background(lipgloss.Color(types.Surface1)).
			Foreground(lipgloss.Color(types.Text)).
			Padding(0, 2)

	NavBarBackground = lipgloss.NewStyle().
				Background(lipgloss.Color(types.Surface1))
)
