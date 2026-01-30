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

type KeychainItem struct {
	ID       uint
	Name     string
	ItemType string // "Key" or "Identity"
	Detail   string
}

type KeychainSidebar struct {
	allItems        []KeychainItem
	filteredItems   []KeychainItem
	selectedIdx     int
	filterActive    bool
	filterText      string
	filterCursorPos int
}

func NewKeychainSidebar() *KeychainSidebar {
	return &KeychainSidebar{
		allItems:        []KeychainItem{},
		filteredItems:   []KeychainItem{},
		selectedIdx:     0,
		filterActive:    false,
		filterText:      "",
		filterCursorPos: 0,
	}
}

func (sidebar *KeychainSidebar) SetItems(keys []models.Key, identities []models.Identity) {
	sidebar.allItems = []KeychainItem{}

	keyIdx := 0
	identityIdx := 0

	for keyIdx < len(keys) || identityIdx < len(identities) {
		if keyIdx >= len(keys) {
			sidebar.allItems = append(sidebar.allItems, KeychainItem{
				ID:       identities[identityIdx].ID,
				Name:     identities[identityIdx].Name,
				ItemType: "Identity",
				Detail:   identities[identityIdx].Username,
			})
			identityIdx++
		} else if identityIdx >= len(identities) {
			sidebar.allItems = append(sidebar.allItems, KeychainItem{
				ID:       keys[keyIdx].ID,
				Name:     keys[keyIdx].Name,
				ItemType: "Key",
				Detail:   "Private Key",
			})
			keyIdx++
		} else if keys[keyIdx].ID > identities[identityIdx].ID {
			sidebar.allItems = append(sidebar.allItems, KeychainItem{
				ID:       keys[keyIdx].ID,
				Name:     keys[keyIdx].Name,
				ItemType: "Key",
				Detail:   "Private Key",
			})
			keyIdx++
		} else {
			sidebar.allItems = append(sidebar.allItems, KeychainItem{
				ID:       identities[identityIdx].ID,
				Name:     identities[identityIdx].Name,
				ItemType: "Identity",
				Detail:   identities[identityIdx].Username,
			})
			identityIdx++
		}
	}

	sidebar.applyFilter()
	if sidebar.selectedIdx >= len(sidebar.filteredItems) {
		sidebar.selectedIdx = 0
	}
}

func (sidebar *KeychainSidebar) IsFilterActive() bool {
	return sidebar.filterActive
}

func (sidebar *KeychainSidebar) SetFilterActive(active bool) {
	sidebar.filterActive = active
}

func (sidebar *KeychainSidebar) GetSelected() *KeychainItem {
	if sidebar.selectedIdx >= 0 && sidebar.selectedIdx < len(sidebar.filteredItems) {
		return &sidebar.filteredItems[sidebar.selectedIdx]
	}
	return nil
}

func (sidebar *KeychainSidebar) SelectItemByID(id uint, itemType string) {
	for i, item := range sidebar.filteredItems {
		if item.ID == id && item.ItemType == itemType {
			sidebar.selectedIdx = i
			return
		}
	}
}

func (sidebar *KeychainSidebar) applyFilter() {
	if sidebar.filterText == "" {
		sidebar.filteredItems = sidebar.allItems
	} else {
		sidebar.filteredItems = []KeychainItem{}
		filterLower := strings.ToLower(sidebar.filterText)
		for _, item := range sidebar.allItems {
			if strings.Contains(strings.ToLower(item.Name), filterLower) {
				sidebar.filteredItems = append(sidebar.filteredItems, item)
			}
		}
	}

	if sidebar.selectedIdx >= len(sidebar.filteredItems) && len(sidebar.filteredItems) > 0 {
		sidebar.selectedIdx = 0
	}
}

func (sidebar *KeychainSidebar) Update(event interface{}) {
	if msg, ok := event.(tea.Msg); ok {
		switch key := msg.(type) {
		case tea.KeyMsg:
			if sidebar.filterActive {
				switch key.Type {
				case tea.KeyEscape:
					var selectedItem *KeychainItem
					if sidebar.selectedIdx >= 0 && sidebar.selectedIdx < len(sidebar.filteredItems) {
						selectedItem = &sidebar.filteredItems[sidebar.selectedIdx]
					}

					sidebar.filterActive = false
					sidebar.filterText = ""
					sidebar.filterCursorPos = 0
					sidebar.applyFilter()

					if selectedItem != nil {
						for i, item := range sidebar.allItems {
							if item.ID == selectedItem.ID && item.ItemType == selectedItem.ItemType {
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
					if sidebar.selectedIdx < len(sidebar.filteredItems)-1 {
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
					if sidebar.selectedIdx < len(sidebar.filteredItems)-1 {
						sidebar.selectedIdx++
					}
				}
			}
		}
	}
}

func (sidebar *KeychainSidebar) Render() string {
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

	if len(sidebar.filteredItems) == 0 {
		noItemsPart := styles.SidebarNormalDesc.Render("No items")
		content = lipgloss.JoinVertical(lipgloss.Left, filterPart, noItemsPart)

		contentHeight := lipgloss.Height(content)
		spacingNeeded := availableHeight - contentHeight + 4
		if spacingNeeded > 0 {
			spacer := lipgloss.NewStyle().Height(spacingNeeded).Render("")
			content = lipgloss.JoinVertical(lipgloss.Left, content, spacer)
		}
	} else {
		totalPages := (len(sidebar.filteredItems) + itemsPerPage - 1) / itemsPerPage
		if totalPages == 0 {
			totalPages = 1
		}

		currentPage := (sidebar.selectedIdx / itemsPerPage)
		pageStartIdx := currentPage * itemsPerPage
		pageEndIdx := pageStartIdx + itemsPerPage
		if pageEndIdx > len(sidebar.filteredItems) {
			pageEndIdx = len(sidebar.filteredItems)
		}

		var itemLines []string
		for i := pageStartIdx; i < pageEndIdx; i++ {
			item := sidebar.filteredItems[i]
			isSelected := i == sidebar.selectedIdx
			line := sidebar.formatItemLine(item, isSelected)
			itemLines = append(itemLines, line)
		}

		if len(itemLines) > 0 {
			itemsContent := lipgloss.JoinVertical(lipgloss.Left, itemLines...)
			content = lipgloss.JoinVertical(lipgloss.Left, filterPart, itemsContent)
		} else {
			content = filterPart
		}

		start := pageStartIdx + 1
		end := pageEndIdx
		totalItems := len(sidebar.filteredItems)
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

func (sidebar *KeychainSidebar) formatItemLine(item KeychainItem, isSelected bool) string {
	title := item.Name
	desc := fmt.Sprintf("[%s] %s", item.ItemType, item.Detail)

	if isSelected {
		styledTitle := styles.SidebarSelectedTitle.Render(title)
		styledDesc := styles.SidebarSelectedDesc.Render(desc)
		itemContent := lipgloss.JoinVertical(lipgloss.Left, styledTitle, styledDesc)
		bordered := styles.SidebarSelectedBorder.Render(itemContent)
		return styles.SidebarItemMargin.Render(bordered)
	}

	styledTitle := styles.SidebarNormalTitle.Render(title)
	styledDesc := styles.SidebarNormalDesc.Render(desc)
	itemContent := lipgloss.JoinVertical(lipgloss.Left, styledTitle, styledDesc)
	padded := styles.SidebarNormalPadding.Render(itemContent)
	return styles.SidebarItemMargin.Render(padded)
}
