package styles

import (
	"yoru/types"

	"github.com/charmbracelet/lipgloss"
)

var (
	TableHeaderCell = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(types.Lavender)).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(lipgloss.Color(types.Surface0)).
			Padding(0, 1)

	TableCell = lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Text)).
			Padding(0, 1)

	TableSelectedRow = lipgloss.NewStyle().
				Background(lipgloss.Color(types.Surface0))

	TableBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(types.Surface0))
)
