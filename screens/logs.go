package screens

import (
	"fmt"
	"strings"
	"time"
	"yoru/repository"
	"yoru/screens/styles"
	"yoru/shared"
	"yoru/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const logsLimit = 50

var logsScreen = &logs{
	selectedIdx: 0,
}

func (screen *logs) Init() tea.Cmd {
	logs, _ := repository.GetLastNConnectionLogs(logsLimit)
	screen.logs = logs
	return nil
}

func (screen *logs) Update(msg tea.Msg) (types.Screen, tea.Cmd) {
	switch message := msg.(type) {
	case tea.KeyMsg:
		switch message.String() {
		case "up":
			if screen.selectedIdx > 0 {
				screen.selectedIdx--
			}
		case "down":
			if screen.selectedIdx < len(screen.logs)-1 {
				screen.selectedIdx++
			}
		}
	}

	return screen, nil
}

func (screen *logs) View() string {
	if len(screen.logs) == 0 {
		emptyMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Subtext0)).
			Render("No connection logs found")
		return lipgloss.Place(
			shared.GlobalState.ScreenWidth,
			shared.GlobalState.ScreenHeight-4,
			lipgloss.Center,
			lipgloss.Center,
			emptyMsg,
		)
	}

	headers := []string{"ID", "Started At", "Ended At", "Local", "Remote", "Mode", "Duration"}
	colWidths := []int{6, 20, 20, 20, 20, 10, 12}

	var headerCells []string
	for i, header := range headers {
		cell := styles.TableHeaderCell.Width(colWidths[i]).Render(header)
		headerCells = append(headerCells, cell)
	}
	headerRow := lipgloss.JoinHorizontal(lipgloss.Top, headerCells...)

	availableHeight := shared.GlobalState.ScreenHeight - 8
	visibleRows := availableHeight
	if visibleRows > len(screen.logs) {
		visibleRows = len(screen.logs)
	}

	startIdx := screen.selectedIdx
	if startIdx+visibleRows > len(screen.logs) {
		startIdx = len(screen.logs) - visibleRows
	}
	if startIdx < 0 {
		startIdx = 0
	}

	var rows []string
	for i := startIdx; i < startIdx+visibleRows && i < len(screen.logs); i++ {
		log := screen.logs[i]
		isSelected := i == screen.selectedIdx

		endedAt := "Active"
		duration := "Ongoing"
		if log.EndedAt != nil {
			endedAt = log.EndedAt.Format("2006-01-02 15:04:05")
			duration = log.EndedAt.Sub(log.StartedAt).Round(time.Second).String()
		}

		cells := []string{
			fmt.Sprintf("%d", log.ID),
			log.StartedAt.Format("2006-01-02 15:04:05"),
			endedAt,
			fmt.Sprintf("%s (%s)", log.LocalHostname, log.LocalIP),
			log.RemoteHostname,
			string(log.Mode),
			duration,
		}

		var rowCells []string
		for j, cell := range cells {
			cellStyle := styles.TableCell.Width(colWidths[j])
			if isSelected {
				cellStyle = cellStyle.Inherit(styles.TableSelectedRow)
			}
			rowCells = append(rowCells, cellStyle.Render(cell))
		}

		row := lipgloss.JoinHorizontal(lipgloss.Top, rowCells...)
		rows = append(rows, row)
	}

	table := lipgloss.JoinVertical(lipgloss.Left, headerRow)
	if len(rows) > 0 {
		table = lipgloss.JoinVertical(lipgloss.Left, headerRow, strings.Join(rows, "\n"))
	}

	bordered := styles.TableBorder.
		Width(shared.GlobalState.ScreenWidth - 4).
		Height(availableHeight).
		Render(table)

	info := lipgloss.NewStyle().
		Foreground(lipgloss.Color(types.Subtext0)).
		Render(fmt.Sprintf("Showing %d of last %d logs | ↑↓: Navigate | g/G: Top/Bottom", len(screen.logs), logsLimit))

	return lipgloss.JoinVertical(lipgloss.Left, bordered, info)
}
