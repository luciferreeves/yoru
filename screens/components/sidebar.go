package components

import (
	"fmt"
	"strings"
	"yoru/models"
	"yoru/screens/styles"
	"yoru/shared"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HostsSidebar struct {
	allHosts        []models.Host
	filteredHosts   []models.Host
	selectedIdx     int
	filterActive    bool
	filterText      string
	filterCursorPos int
}

func NewHostsSidebar() *HostsSidebar {
	return &HostsSidebar{
		allHosts:        []models.Host{},
		filteredHosts:   []models.Host{},
		selectedIdx:     0,
		filterActive:    false,
		filterText:      "",
		filterCursorPos: 0,
	}
}

func (sidebar *HostsSidebar) SetHosts(hosts []models.Host) {
	sidebar.allHosts = hosts
	sidebar.applyFilter()
	if sidebar.selectedIdx >= len(sidebar.filteredHosts) {
		sidebar.selectedIdx = 0
	}
}

func (sidebar *HostsSidebar) IsFilterActive() bool {
	return sidebar.filterActive
}

func (sidebar *HostsSidebar) SetFilterActive(active bool) {
	sidebar.filterActive = active
}

func (sidebar *HostsSidebar) GetSelected() *models.Host {
	if sidebar.selectedIdx >= 0 && sidebar.selectedIdx < len(sidebar.filteredHosts) {
		return &sidebar.filteredHosts[sidebar.selectedIdx]
	}
	return nil
}

func (sidebar *HostsSidebar) applyFilter() {
	if sidebar.filterText == "" {
		sidebar.filteredHosts = sidebar.allHosts
	} else {
		sidebar.filteredHosts = []models.Host{}
		filterLower := strings.ToLower(sidebar.filterText)
		for _, host := range sidebar.allHosts {
			if strings.Contains(strings.ToLower(host.Name), filterLower) {
				sidebar.filteredHosts = append(sidebar.filteredHosts, host)
			}
		}
	}

	if sidebar.selectedIdx >= len(sidebar.filteredHosts) && len(sidebar.filteredHosts) > 0 {
		sidebar.selectedIdx = 0
	}
}

func (sidebar *HostsSidebar) Update(event interface{}) {
	if msg, ok := event.(tea.Msg); ok {
		switch key := msg.(type) {
		case tea.KeyMsg:
			if sidebar.filterActive {
				switch key.Type {
				case tea.KeyEscape:
					var selectedHost *models.Host
					if sidebar.selectedIdx >= 0 && sidebar.selectedIdx < len(sidebar.filteredHosts) {
						selectedHost = &sidebar.filteredHosts[sidebar.selectedIdx]
					}

					sidebar.filterActive = false
					sidebar.filterText = ""
					sidebar.filterCursorPos = 0
					sidebar.applyFilter()

					if selectedHost != nil {
						for i, host := range sidebar.allHosts {
							if host.ID == selectedHost.ID {
								sidebar.selectedIdx = i
								break
							}
						}
					}
				case tea.KeyBackspace:
					if sidebar.filterCursorPos > 0 {
						sidebar.filterText = sidebar.filterText[:sidebar.filterCursorPos-1] + sidebar.filterText[sidebar.filterCursorPos:]
						sidebar.filterCursorPos--
						sidebar.applyFilter()
						sidebar.selectedIdx = 0
					}
				case tea.KeyDelete:
					if sidebar.filterCursorPos < len(sidebar.filterText) {
						sidebar.filterText = sidebar.filterText[:sidebar.filterCursorPos] + sidebar.filterText[sidebar.filterCursorPos+1:]
						sidebar.applyFilter()
						sidebar.selectedIdx = 0
					}
				case tea.KeyLeft:
					if sidebar.filterCursorPos > 0 {
						sidebar.filterCursorPos--
					}
				case tea.KeyRight:
					if sidebar.filterCursorPos < len(sidebar.filterText) {
						sidebar.filterCursorPos++
					}
				case tea.KeyHome:
					sidebar.filterCursorPos = 0
				case tea.KeyEnd:
					sidebar.filterCursorPos = len(sidebar.filterText)
				case tea.KeyUp:
					if sidebar.selectedIdx > 0 {
						sidebar.selectedIdx--
					}
				case tea.KeyDown:
					if sidebar.selectedIdx < len(sidebar.filteredHosts)-1 {
						sidebar.selectedIdx++
					}
				default:
					if len(key.Runes) > 0 && key.Runes[0] >= 32 && key.Runes[0] < 127 {
						sidebar.filterText = sidebar.filterText[:sidebar.filterCursorPos] + string(key.Runes[0]) + sidebar.filterText[sidebar.filterCursorPos:]
						sidebar.filterCursorPos++
						sidebar.applyFilter()
						sidebar.selectedIdx = 0
					}
				}
			} else {
				itemsPerPage := (shared.GlobalState.ScreenHeight - 5) / 3
				if itemsPerPage < 1 {
					itemsPerPage = 1
				}

				switch key.String() {
				case "/":
					sidebar.filterActive = true
					sidebar.filterText = ""
					sidebar.filterCursorPos = 0
				case "up":
					if sidebar.selectedIdx > 0 {
						sidebar.selectedIdx--
					}
				case "down":
					if sidebar.selectedIdx < len(sidebar.filteredHosts)-1 {
						sidebar.selectedIdx++
					}
				}
			}
		}
	}
}

func (sidebar *HostsSidebar) Render() string {
	availableHeight := shared.GlobalState.ScreenHeight - 8

	itemsPerPage := max((availableHeight)/3, 1)

	var filterPart string
	if sidebar.filterActive {
		before := sidebar.filterText[:sidebar.filterCursorPos]
		after := ""
		if sidebar.filterCursorPos < len(sidebar.filterText) {
			after = sidebar.filterText[sidebar.filterCursorPos+1:]
		}
		cursorChar := " "
		if sidebar.filterCursorPos < len(sidebar.filterText) {
			cursorChar = string(sidebar.filterText[sidebar.filterCursorPos])
		}
		cursor := styles.SidebarFilterCursor.Render(cursorChar)
		filterDisplay := before + cursor + after
		filterRendered := styles.SidebarFilterActive.Render("/ " + filterDisplay)
		filterPart = styles.SidebarFilterWithMargin.Render(filterRendered)
	} else {
		filterRendered := styles.SidebarFilterInactive.Render("/ ")
		filterPart = styles.SidebarFilterWithMargin.Render(filterRendered)
	}

	var content string
	var bottomContent string

	if len(sidebar.filteredHosts) == 0 {
		noHostsPart := styles.SidebarNormalDesc.Render("No hosts")
		content = lipgloss.JoinVertical(lipgloss.Left, filterPart, noHostsPart)
	} else {
		totalPages := (len(sidebar.filteredHosts) + itemsPerPage - 1) / itemsPerPage
		if totalPages == 0 {
			totalPages = 1
		}

		currentPage := (sidebar.selectedIdx / itemsPerPage)
		pageStartIdx := currentPage * itemsPerPage
		pageEndIdx := pageStartIdx + itemsPerPage
		if pageEndIdx > len(sidebar.filteredHosts) {
			pageEndIdx = len(sidebar.filteredHosts)
		}

		var hostItems []string
		for i := pageStartIdx; i < pageEndIdx; i++ {
			host := sidebar.filteredHosts[i]
			isSelected := i == sidebar.selectedIdx
			line := sidebar.formatHostLine(host, isSelected)
			hostItems = append(hostItems, line)
		}

		if len(hostItems) > 0 {
			itemsContent := lipgloss.JoinVertical(lipgloss.Left, hostItems...)
			content = lipgloss.JoinVertical(lipgloss.Left, filterPart, itemsContent)
		} else {
			content = filterPart
		}

		start := pageStartIdx + 1
		end := pageEndIdx
		totalItems := len(sidebar.filteredHosts)
		pageNum := currentPage + 1

		var dots strings.Builder
		for p := 0; p < totalPages; p++ {
			if p == currentPage {
				dots.WriteString(styles.SidebarPaginationActive.Render("●"))
			} else {
				dots.WriteString(styles.SidebarPaginationInactive.Render("○"))
			}
			if p < totalPages-1 {
				dots.WriteString(" ")
			}
		}

		paginationText := styles.SidebarPaginationText.Render(
			fmt.Sprintf("%d-%d of %d (%d/%d)", start, end, totalItems, pageNum, totalPages))

		bottomContent = lipgloss.JoinVertical(lipgloss.Left, dots.String(), paginationText)
	}

	if bottomContent != "" {
		contentHeight := lipgloss.Height(content)
		bottomHeight := lipgloss.Height(bottomContent)

		totalUsedHeight := contentHeight + bottomHeight
		spacingNeeded := availableHeight - totalUsedHeight + 4

		if spacingNeeded > 0 {
			spacingStyle := lipgloss.NewStyle().MarginTop(spacingNeeded)
			bottomWithMargin := spacingStyle.Render(bottomContent)
			return lipgloss.JoinVertical(lipgloss.Left, content, bottomWithMargin)
		}
		return lipgloss.JoinVertical(lipgloss.Left, content, bottomContent)
	}

	return content
}

func (sidebar *HostsSidebar) formatHostLine(host models.Host, isSelected bool) string {
	title := host.Name
	desc := ""
	if host.Hostname == "" {
		desc = fmt.Sprintf(":%d", host.Port)
	} else {
		desc = fmt.Sprintf("%s:%d", host.Hostname, host.Port)
	}

	if isSelected {
		styledTitle := styles.SidebarSelectedTitle.Render(title)
		styledDesc := styles.SidebarSelectedDesc.Render(desc)
		item := lipgloss.JoinVertical(lipgloss.Left, styledTitle, styledDesc)
		bordered := styles.SidebarSelectedBorder.Render(item)
		return styles.SidebarItemMargin.Render(bordered)
	}

	styledTitle := styles.SidebarNormalTitle.Render(title)
	styledDesc := styles.SidebarNormalDesc.Render(desc)
	item := lipgloss.JoinVertical(lipgloss.Left, styledTitle, styledDesc)
	padded := styles.SidebarNormalPadding.Render(item)
	return styles.SidebarItemMargin.Render(padded)
}
