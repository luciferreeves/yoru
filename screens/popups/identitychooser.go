package popups

import (
	"strings"
	"yoru/repository"
	"yoru/screens/components"
	"yoru/screens/styles"
	"yoru/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CredentialItem struct {
	ID      uint
	Name    string
	Details string
	Type    types.CredentialType
}

type IdentityChooserPopup struct {
	popup         *components.Popup
	items         []CredentialItem
	filteredItems []CredentialItem
	selectedIdx   int
	filterActive  bool
	filterText    string
	onSelect      func(types.CredentialType, uint)
	onCancel      func()
}

func NewIdentityChooserPopup() *IdentityChooserPopup {
	icp := &IdentityChooserPopup{
		popup: components.NewPopup(),
	}
	return icp
}

func (icp *IdentityChooserPopup) Show(currentCredentialType types.CredentialType, currentCredentialID uint, onSelect func(types.CredentialType, uint), onCancel func()) {
	// Load both identities and keys
	icp.items = []CredentialItem{}

	identities, _ := repository.GetAllIdentities()
	for _, identity := range identities {
		icp.items = append(icp.items, CredentialItem{
			ID:      identity.ID,
			Name:    identity.Name,
			Details: identity.Username,
			Type:    types.CredentialIdentity,
		})
	}

	keys, _ := repository.GetAllKeys()
	for _, key := range keys {
		icp.items = append(icp.items, CredentialItem{
			ID:      key.ID,
			Name:    key.Name,
			Details: "SSH Key",
			Type:    types.CredentialKey,
		})
	}

	icp.filteredItems = icp.items
	icp.onSelect = onSelect
	icp.onCancel = onCancel
	icp.filterActive = false
	icp.filterText = ""
	icp.selectedIdx = 0

	// Find the currently selected credential
	for i, item := range icp.filteredItems {
		if item.Type == currentCredentialType && item.ID == currentCredentialID {
			icp.selectedIdx = i + 1
			break
		}
	}

	icp.popup.Show(icp.buildContent(), icp.handleInput)
}

func (icp *IdentityChooserPopup) Hide() {
	icp.popup.Hide()
}

func (icp *IdentityChooserPopup) IsVisible() bool {
	return icp.popup.IsVisible()
}

func (icp *IdentityChooserPopup) Update(msg tea.Msg) {
	icp.popup.Update(msg)
}

func (icp *IdentityChooserPopup) Render() string {
	return icp.popup.Render()
}

func (icp *IdentityChooserPopup) applyFilter() {
	if icp.filterText == "" {
		icp.filteredItems = icp.items
	} else {
		icp.filteredItems = []CredentialItem{}
		filterLower := strings.ToLower(icp.filterText)
		for _, item := range icp.items {
			if strings.Contains(strings.ToLower(item.Name), filterLower) ||
				strings.Contains(strings.ToLower(item.Details), filterLower) {
				icp.filteredItems = append(icp.filteredItems, item)
			}
		}
	}

	if icp.selectedIdx > len(icp.filteredItems) {
		icp.selectedIdx = 0
	}
}

func (icp *IdentityChooserPopup) handleInput(msg tea.Msg) bool {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if icp.filterActive {
			switch keyMsg.Type {
			case tea.KeyEscape:
				icp.filterActive = false
				icp.filterText = ""
				icp.applyFilter()
				icp.popup.SetContent(icp.buildContent())
				return true
			case tea.KeyBackspace:
				if len(icp.filterText) > 0 {
					icp.filterText = icp.filterText[:len(icp.filterText)-1]
					icp.applyFilter()
					icp.popup.SetContent(icp.buildContent())
				}
				return true
			case tea.KeyRunes:
				icp.filterText += string(keyMsg.Runes)
				icp.applyFilter()
				icp.popup.SetContent(icp.buildContent())
				return true
			case tea.KeyUp:
				if icp.selectedIdx > 0 {
					icp.selectedIdx--
					icp.popup.SetContent(icp.buildContent())
				}
				return true
			case tea.KeyDown:
				if icp.selectedIdx < len(icp.filteredItems) {
					icp.selectedIdx++
					icp.popup.SetContent(icp.buildContent())
				}
				return true
			case tea.KeyEnter:
				if icp.selectedIdx == 0 {
					if icp.onSelect != nil {
						icp.onSelect("", 0)
					}
				} else if icp.selectedIdx > 0 && icp.selectedIdx <= len(icp.filteredItems) {
					item := icp.filteredItems[icp.selectedIdx-1]
					if icp.onSelect != nil {
						icp.onSelect(item.Type, item.ID)
					}
				}
				icp.Hide()
				return true
			}
		} else {
			switch keyMsg.String() {
			case "/":
				icp.filterActive = true
				icp.popup.SetContent(icp.buildContent())
				return true
			case "up":
				if icp.selectedIdx > 0 {
					icp.selectedIdx--
					icp.popup.SetContent(icp.buildContent())
				}
				return true
			case "down":
				if icp.selectedIdx < len(icp.filteredItems) {
					icp.selectedIdx++
					icp.popup.SetContent(icp.buildContent())
				}
				return true
			case "enter":
				if icp.selectedIdx == 0 {
					if icp.onSelect != nil {
						icp.onSelect("", 0)
					}
				} else if icp.selectedIdx > 0 && icp.selectedIdx <= len(icp.filteredItems) {
					item := icp.filteredItems[icp.selectedIdx-1]
					if icp.onSelect != nil {
						icp.onSelect(item.Type, item.ID)
					}
				}
				icp.Hide()
				return true
			case "esc":
				if icp.onCancel != nil {
					icp.onCancel()
				}
				icp.Hide()
				return true
			}
		}
	}
	return false
}

func (icp *IdentityChooserPopup) buildContent() string {
	title := styles.PopupTitle.Render("Select Identity")

	var filterLine string
	if icp.filterActive {
		filterLine = styles.SidebarFilterActive.Render("/ " + icp.filterText + "█")
	} else {
		filterLine = styles.SidebarFilterInactive.Render("/ ")
	}

	var items []string

	// Clear Identity option
	if icp.selectedIdx == 0 {
		items = append(items, styles.PopupItemSelected.Render("Clear Identity"))
	} else {
		items = append(items, styles.PopupItemNormal.Render("Clear Identity"))
	}

	// Show all items without pagination/scrolling
	for i, item := range icp.filteredItems {
		displayText := item.Name
		if item.Details != "" {
			displayText += " (" + item.Details + ")"
		}

		if i+1 == icp.selectedIdx {
			items = append(items, styles.PopupItemSelected.Render(displayText))
		} else {
			items = append(items, styles.PopupItemNormal.Render(displayText))
		}
	}

	if len(icp.filteredItems) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Subtext0))
		items = append(items, emptyStyle.Render("No credentials found"))
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		filterLine,
		"",
		lipgloss.JoinVertical(lipgloss.Left, items...),
	)

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(content)
}
