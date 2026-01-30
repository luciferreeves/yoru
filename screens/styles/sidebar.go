package styles

import (
	"yoru/types"

	"github.com/charmbracelet/lipgloss"
)

var (
	SidebarSelectedTitle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(types.Lavender)).
				Bold(true)

	SidebarSelectedDesc = lipgloss.NewStyle().
				Foreground(lipgloss.Color(types.Sky))
	SidebarNormalTitle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(types.Text))

	SidebarNormalDesc = lipgloss.NewStyle().
				Foreground(lipgloss.Color(types.Subtext0))
	SidebarSelectedBorder = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderLeft(true).
				BorderForeground(lipgloss.Color(types.Lavender)).
				PaddingLeft(1)

	SidebarNormalPadding = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderLeft(true).
				BorderForeground(lipgloss.Color(types.Base)).
				PaddingLeft(1)

	SidebarItemMargin = lipgloss.NewStyle().
				MarginBottom(1)

	SidebarFilterActive = lipgloss.NewStyle().
				Foreground(lipgloss.Color(types.Lavender)).
				Bold(true)

	SidebarFilterInactive = lipgloss.NewStyle().
				Foreground(lipgloss.Color(types.Subtext0))

	SidebarFilterWithMargin = lipgloss.NewStyle().
				MarginBottom(1)

	SidebarFilterCursor = lipgloss.NewStyle().
				Background(lipgloss.Color(types.Lavender)).
				Foreground(lipgloss.Color(types.Base))

	SidebarPaginationActive = lipgloss.NewStyle().
				Foreground(lipgloss.Color(types.Lavender))

	SidebarPaginationInactive = lipgloss.NewStyle().
					Foreground(lipgloss.Color(types.Surface1))

	SidebarPaginationText = lipgloss.NewStyle().
				Foreground(lipgloss.Color(types.Subtext0))
)
