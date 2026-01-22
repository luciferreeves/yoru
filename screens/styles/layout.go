package styles

import (
	"yoru/types"

	"github.com/charmbracelet/lipgloss"
)

var (
	ContentBackground = lipgloss.NewStyle().
				Background(lipgloss.Color(types.Base)).
				Foreground(lipgloss.Color(types.Text))

	ContentArea = lipgloss.NewStyle().
			Background(lipgloss.Color(types.Base)).
			Foreground(lipgloss.Color(types.Text)).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(types.Surface2)).
			Padding(0, 1)

	SidebarArea = lipgloss.NewStyle().
			Background(lipgloss.Color(types.Base)).
			BorderStyle(lipgloss.NormalBorder()).
			BorderRight(true).
			BorderForeground(lipgloss.Color(types.Surface2))

	FormArea = lipgloss.NewStyle().
			Background(lipgloss.Color(types.Base))
)
