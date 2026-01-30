package components

import (
	"strings"
	"yoru/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TextArea struct {
	value       string
	placeholder string
	width       int
	height      int
	isEditing   bool
	isFocused   bool

	cursor int

	viewportLine          int
	viewportWrappedOffset int
}

func NewTextArea(placeholder string, width, height int) TextArea {
	return TextArea{
		value:                 "",
		placeholder:           placeholder,
		width:                 width,
		height:                height,
		isEditing:             false,
		isFocused:             false,
		cursor:                0,
		viewportLine:          0,
		viewportWrappedOffset: 0,
	}
}

func (ta *TextArea) Focus() {
	ta.isFocused = true
}

func (ta *TextArea) Blur() {
	ta.isFocused = false
	ta.isEditing = false
}

func (ta *TextArea) StartEditing() {
	ta.isEditing = true
	ta.cursor = len(ta.value)
	ta.scrollToEnd()
	ta.adjustViewport()
}

func (ta *TextArea) StopEditing() {
	ta.isEditing = false
	ta.cursor = 0
	ta.viewportLine = 0
	ta.viewportWrappedOffset = 0
}

func (ta *TextArea) IsEditing() bool {
	return ta.isEditing
}

func (ta *TextArea) Value() string {
	return ta.value
}

func (ta *TextArea) SetValue(value string) {
	ta.value = value
	if ta.cursor > len(ta.value) {
		ta.cursor = len(ta.value)
	}
}

func (ta *TextArea) Reset() {
	ta.value = ""
	ta.cursor = 0
	ta.viewportLine = 0
	ta.viewportWrappedOffset = 0
	ta.isEditing = false
}

func (ta *TextArea) Update(msg tea.Msg) tea.Cmd {
	if !ta.isEditing {
		return nil
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}

	if keyMsg.Alt && keyMsg.Type == tea.KeyLeft {
		ta.moveCursorWordLeft()
		return nil
	}
	if keyMsg.Alt && keyMsg.Type == tea.KeyRight {
		ta.moveCursorWordRight()
		return nil
	}

	if keyMsg.Alt && keyMsg.Type == tea.KeyBackspace {
		ta.deleteWordBackward()
		return nil
	}
	if keyMsg.Alt && keyMsg.Type == tea.KeyDelete {
		ta.deleteWordForward()
		return nil
	}

	switch keyMsg.Type {
	case tea.KeyBackspace:
		if ta.cursor > 0 {
			ta.value = ta.value[:ta.cursor-1] + ta.value[ta.cursor:]
			ta.cursor--
			ta.adjustViewport()
		}
	case tea.KeyDelete:
		if ta.cursor < len(ta.value) {
			ta.value = ta.value[:ta.cursor] + ta.value[ta.cursor+1:]
			ta.adjustViewport()
		}
	case tea.KeyLeft:
		if ta.cursor > 0 {
			ta.cursor--
			ta.adjustViewport()
		}
	case tea.KeyRight:
		if ta.cursor < len(ta.value) {
			ta.cursor++
			ta.adjustViewport()
		}
	case tea.KeyUp:
		ta.moveCursorUp()
	case tea.KeyDown:
		ta.moveCursorDown()
	case tea.KeyHome, tea.KeyCtrlA:
		ta.moveCursorToLineStart()
	case tea.KeyEnd, tea.KeyCtrlE:
		ta.moveCursorToLineEnd()
	case tea.KeyEnter:
		ta.value = ta.value[:ta.cursor] + "\n" + ta.value[ta.cursor:]
		ta.cursor++
		ta.adjustViewport()
	case tea.KeyRunes:
		insertedText := string(keyMsg.Runes)
		ta.value = ta.value[:ta.cursor] + insertedText + ta.value[ta.cursor:]
		ta.cursor += len(insertedText)
		ta.adjustViewport()
	case tea.KeySpace:
		ta.value = ta.value[:ta.cursor] + " " + ta.value[ta.cursor:]
		ta.cursor++
	case tea.KeyTab:
		ta.value = ta.value[:ta.cursor] + "    " + ta.value[ta.cursor:]
		ta.cursor += 4
	}

	return nil
}

func (ta *TextArea) View() string {
	lines := ta.getLines()

	if len(ta.value) == 0 {
		placeholderLines := make([]string, ta.height)
		placeholderLines[0] = ta.placeholder
		return ta.renderFixedSize(placeholderLines, true)
	}

	visibleLines := ta.getVisibleLines(lines)

	if ta.isEditing {
		visibleLines = ta.addCursor(visibleLines)
	}

	return ta.renderFixedSize(visibleLines, false)
}

func (ta *TextArea) renderFixedSize(lines []string, isPlaceholder bool) string {
	for i := range lines {
		displayWidth := lipgloss.Width(lines[i])

		if displayWidth < ta.width {
			lines[i] += strings.Repeat(" ", ta.width-displayWidth)
		} else if displayWidth > ta.width {
			runes := []rune(lines[i])
			if len(runes) > ta.width {
				lines[i] = string(runes[:ta.width])
			}
		}
	}

	content := strings.Join(lines, "\n")

	if isPlaceholder {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(types.Overlay0)).
			Render(content)
	}

	return content
}

func (ta *TextArea) wrapLine(line string) []string {
	effectiveWidth := ta.width - 1
	if effectiveWidth < 1 {
		effectiveWidth = 1
	}

	if lipgloss.Width(line) <= effectiveWidth {
		return []string{line}
	}

	displayWidth := lipgloss.Width(line)
	if displayWidth != len(line) {
		runes := []rune(line)
		if len(runes) > effectiveWidth {
			return []string{string(runes[:effectiveWidth])}
		}
		return []string{line}
	}

	result := []string{}
	runes := []rune(line)

	for len(runes) > 0 {
		if len(runes) <= effectiveWidth {
			result = append(result, string(runes))
			break
		}
		result = append(result, string(runes[:effectiveWidth]))
		runes = runes[effectiveWidth:]
	}

	return result
}

func (ta *TextArea) getLines() []string {
	if len(ta.value) == 0 {
		emptyLines := make([]string, ta.height)
		return emptyLines
	}
	return strings.Split(ta.value, "\n")
}

func (ta *TextArea) getVisibleLines(lines []string) []string {
	if len(lines) == 0 {
		emptyLines := make([]string, ta.height)
		return emptyLines
	}

	start := ta.viewportLine
	if start >= len(lines) {
		start = 0
	}

	wrapped := []string{}
	for i := start; i < len(lines) && len(wrapped) < ta.height; i++ {
		wrappedLine := ta.wrapLine(lines[i])
		startIdx := 0
		if i == start {
			startIdx = ta.viewportWrappedOffset
			if startIdx >= len(wrappedLine) {
				startIdx = 0
			}
		}
		for j := startIdx; j < len(wrappedLine); j++ {
			if len(wrapped) < ta.height {
				wrapped = append(wrapped, wrappedLine[j])
			} else {
				break
			}
		}
	}

	for len(wrapped) < ta.height {
		wrapped = append(wrapped, "")
	}

	return wrapped
}

func (ta *TextArea) addCursor(wrappedLines []string) []string {
	lines := ta.getLines()
	cursorLine, cursorCol := ta.getCursorPosition()

	displayLineIdx := 0
	remainingCol := cursorCol

	for i := 0; i < ta.viewportLine && i < len(lines); i++ {
		wrapped := ta.wrapLine(lines[i])
		displayLineIdx += len(wrapped)
	}

	for i := ta.viewportLine; i < cursorLine && i < len(lines); i++ {
		wrapped := ta.wrapLine(lines[i])
		displayLineIdx += len(wrapped)
	}

	if cursorLine < len(lines) {
		wrapped := ta.wrapLine(lines[cursorLine])
		for wIdx, wrappedSegment := range wrapped {
			segmentLen := len([]rune(wrappedSegment))
			if remainingCol <= segmentLen || wIdx == len(wrapped)-1 {
				wrappedCount := 0

				if ta.viewportLine == cursorLine {
					wrappedCount = wIdx - ta.viewportWrappedOffset
				} else {
					if ta.viewportLine < len(lines) {
						firstWrapped := ta.wrapLine(lines[ta.viewportLine])
						wrappedCount -= ta.viewportWrappedOffset
						wrappedCount += len(firstWrapped)
					}
					for i := ta.viewportLine + 1; i < cursorLine && i < len(lines); i++ {
						w := ta.wrapLine(lines[i])
						wrappedCount += len(w)
					}
					wrappedCount += wIdx
				}

				if wrappedCount >= 0 && wrappedCount < len(wrappedLines) {
					line := wrappedLines[wrappedCount]
					cursorStyle := lipgloss.NewStyle().Background(lipgloss.Color(types.Text)).Foreground(lipgloss.Color(types.Base))
					if remainingCol >= len([]rune(line)) {
						wrappedLines[wrappedCount] = line + cursorStyle.Render(" ")
					} else {
						runes := []rune(line)
						if remainingCol < len(runes) {
							wrappedLines[wrappedCount] = string(runes[:remainingCol]) + cursorStyle.Render(string(runes[remainingCol])) + string(runes[remainingCol+1:])
						}
					}
				}
				break
			}
			remainingCol -= segmentLen
		}
	}

	return wrappedLines
}

func (ta *TextArea) getCursorPosition() (line, col int) {
	line = 0
	col = 0

	for i := 0; i < ta.cursor && i < len(ta.value); i++ {
		if ta.value[i] == '\n' {
			line++
			col = 0
		} else {
			col++
		}
	}

	return line, col
}

func (ta *TextArea) moveCursorUp() {
	line, col := ta.getCursorPosition()
	lines := ta.getLines()

	if line == 0 && col == 0 {
		return
	}

	effectiveWidth := ta.width - 1
	if effectiveWidth < 1 {
		effectiveWidth = 1
	}

	var displayLines []struct {
		lineIndex    int
		wrappedIndex int
		startCol     int
		endCol       int
	}

	for i := 0; i <= line; i++ {
		wrapped := ta.wrapLine(lines[i])
		for w := 0; w < len(wrapped); w++ {
			startCol := w * effectiveWidth
			endCol := startCol + len(wrapped[w])
			if endCol > len(lines[i]) {
				endCol = len(lines[i])
			}
			displayLines = append(displayLines, struct {
				lineIndex    int
				wrappedIndex int
				startCol     int
				endCol       int
			}{i, w, startCol, endCol})
		}
	}

	currentDisplayLine := -1
	for i, dl := range displayLines {
		if dl.lineIndex == line && col >= dl.startCol && col <= dl.endCol {
			currentDisplayLine = i
			break
		}
	}

	if currentDisplayLine <= 0 {
		ta.cursor = 0
		ta.adjustViewport()
		return
	}

	prevDisplayLine := displayLines[currentDisplayLine-1]
	currentDisplayInfo := displayLines[currentDisplayLine]

	colInCurrentDisplay := col - currentDisplayInfo.startCol
	targetCol := prevDisplayLine.startCol + colInCurrentDisplay

	if targetCol > prevDisplayLine.endCol {
		targetCol = prevDisplayLine.endCol
	}
	if targetCol < prevDisplayLine.startCol {
		targetCol = prevDisplayLine.startCol
	}

	ta.cursor = ta.getOffsetFromPosition(prevDisplayLine.lineIndex, targetCol)
	ta.adjustViewport()
}

func (ta *TextArea) moveCursorDown() {
	line, col := ta.getCursorPosition()
	lines := ta.getLines()

	if line >= len(lines)-1 && col >= len(lines[line]) {
		return
	}

	effectiveWidth := ta.width - 1
	if effectiveWidth < 1 {
		effectiveWidth = 1
	}

	var displayLines []struct {
		lineIndex    int
		wrappedIndex int
		startCol     int
		endCol       int
	}

	for i := 0; i < len(lines); i++ {
		wrapped := ta.wrapLine(lines[i])
		for w := 0; w < len(wrapped); w++ {
			startCol := w * effectiveWidth
			endCol := startCol + len(wrapped[w])
			if endCol > len(lines[i]) {
				endCol = len(lines[i])
			}
			displayLines = append(displayLines, struct {
				lineIndex    int
				wrappedIndex int
				startCol     int
				endCol       int
			}{i, w, startCol, endCol})
		}
	}

	currentDisplayLine := -1
	for i, dl := range displayLines {
		if dl.lineIndex == line && col >= dl.startCol && col <= dl.endCol {
			currentDisplayLine = i
			break
		}
	}

	if currentDisplayLine < 0 || currentDisplayLine >= len(displayLines)-1 {
		ta.cursor = len(ta.value)
		ta.adjustViewport()
		return
	}

	nextDisplayLine := displayLines[currentDisplayLine+1]
	currentDisplayInfo := displayLines[currentDisplayLine]

	colInCurrentDisplay := col - currentDisplayInfo.startCol
	targetCol := nextDisplayLine.startCol + colInCurrentDisplay

	if targetCol > nextDisplayLine.endCol {
		targetCol = nextDisplayLine.endCol
	}
	if targetCol < nextDisplayLine.startCol {
		targetCol = nextDisplayLine.startCol
	}

	ta.cursor = ta.getOffsetFromPosition(nextDisplayLine.lineIndex, targetCol)
	ta.adjustViewport()
}

func (ta *TextArea) moveCursorToLineStart() {
	line, _ := ta.getCursorPosition()
	ta.cursor = ta.getOffsetFromPosition(line, 0)
	ta.adjustViewport()
}

func (ta *TextArea) moveCursorToLineEnd() {
	line, _ := ta.getCursorPosition()
	lines := ta.getLines()
	if line < len(lines) {
		ta.cursor = ta.getOffsetFromPosition(line, len(lines[line]))
	}
	ta.adjustViewport()
}

func (ta *TextArea) getOffsetFromPosition(line, col int) int {
	offset := 0
	lines := ta.getLines()

	for i := 0; i < line && i < len(lines); i++ {
		offset += len(lines[i]) + 1
	}

	if line < len(lines) && col <= len(lines[line]) {
		offset += col
	}

	return offset
}

func (ta *TextArea) moveCursorWordLeft() {
	if ta.cursor == 0 {
		return
	}

	for ta.cursor > 0 && !isWordChar(rune(ta.value[ta.cursor-1])) {
		ta.cursor--
	}

	for ta.cursor > 0 && isWordChar(rune(ta.value[ta.cursor-1])) {
		ta.cursor--
	}

	ta.adjustViewport()
}

func (ta *TextArea) moveCursorWordRight() {
	if ta.cursor >= len(ta.value) {
		return
	}

	for ta.cursor < len(ta.value) && !isWordChar(rune(ta.value[ta.cursor])) {
		ta.cursor++
	}

	for ta.cursor < len(ta.value) && isWordChar(rune(ta.value[ta.cursor])) {
		ta.cursor++
	}

	ta.adjustViewport()
}

func (ta *TextArea) deleteWordBackward() {
	if ta.cursor == 0 {
		return
	}

	startCursor := ta.cursor

	for ta.cursor > 0 && !isWordChar(rune(ta.value[ta.cursor-1])) {
		ta.cursor--
	}

	for ta.cursor > 0 && isWordChar(rune(ta.value[ta.cursor-1])) {
		ta.cursor--
	}

	ta.value = ta.value[:ta.cursor] + ta.value[startCursor:]
	ta.adjustViewport()
}

func (ta *TextArea) deleteWordForward() {
	if ta.cursor >= len(ta.value) {
		return
	}

	startCursor := ta.cursor

	for ta.cursor < len(ta.value) && !isWordChar(rune(ta.value[ta.cursor])) {
		ta.cursor++
	}

	for ta.cursor < len(ta.value) && isWordChar(rune(ta.value[ta.cursor])) {
		ta.cursor++
	}

	ta.value = ta.value[:startCursor] + ta.value[ta.cursor:]
	ta.cursor = startCursor
	ta.adjustViewport()
}

func isWordChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
}

func (ta *TextArea) adjustViewport() {
	line, col := ta.getCursorPosition()
	lines := ta.getLines()
	totalLines := len(lines)

	if totalLines == 0 || line >= totalLines {
		ta.viewportLine = 0
		return
	}

	displayLine := 0
	for i := 0; i < line && i < len(lines); i++ {
		wrappedLines := ta.wrapLine(lines[i])
		displayLine += len(wrappedLines)
	}

	if line < len(lines) {
		wrappedLines := ta.wrapLine(lines[line])
		effectiveWidth := ta.width - 1
		if effectiveWidth < 1 {
			effectiveWidth = 1
		}
		wrappedSegmentIndex := col / effectiveWidth
		if wrappedSegmentIndex >= len(wrappedLines) {
			wrappedSegmentIndex = len(wrappedLines) - 1
		}
		displayLine += wrappedSegmentIndex
	}

	viewportDisplayLines := 0
	for i := 0; i < ta.viewportLine && i < len(lines); i++ {
		wrappedLines := ta.wrapLine(lines[i])
		viewportDisplayLines += len(wrappedLines)
	}

	viewportDisplayLines += ta.viewportWrappedOffset

	if displayLine >= viewportDisplayLines+ta.height {
		targetDisplayLine := displayLine - ta.height + 1
		if targetDisplayLine < 0 {
			targetDisplayLine = 0
		}
		accumulatedLines := 0
		for i := 0; i < len(lines); i++ {
			wrappedLines := ta.wrapLine(lines[i])

			if accumulatedLines+len(wrappedLines) > targetDisplayLine {
				ta.viewportLine = i

				ta.viewportWrappedOffset = targetDisplayLine - accumulatedLines
				if ta.viewportWrappedOffset < 0 {
					ta.viewportWrappedOffset = 0
				}
				break
			}
			accumulatedLines += len(wrappedLines)
			if i == len(lines)-1 {
				ta.viewportLine = i
				ta.viewportWrappedOffset = 0
			}
		}
	}

	if displayLine < viewportDisplayLines {
		accumulated := 0
		for i := 0; i < len(lines); i++ {
			wrapped := ta.wrapLine(lines[i])
			if accumulated+len(wrapped) > displayLine {
				ta.viewportLine = i

				ta.viewportWrappedOffset = displayLine - accumulated
				if ta.viewportWrappedOffset < 0 {
					ta.viewportWrappedOffset = 0
				}
				break
			}
			accumulated += len(wrapped)
		}
	}

	displayLinesFromViewport := 0
	if ta.viewportLine < len(lines) {
		firstWrapped := ta.wrapLine(lines[ta.viewportLine])

		if ta.viewportWrappedOffset < len(firstWrapped) {
			displayLinesFromViewport += len(firstWrapped) - ta.viewportWrappedOffset
		}
	}
	for i := ta.viewportLine + 1; i < len(lines); i++ {
		wrapped := ta.wrapLine(lines[i])
		displayLinesFromViewport += len(wrapped)
	}

	if displayLinesFromViewport < ta.height && totalLines > 0 {
		targetDisplayLines := 0
		newViewport := len(lines) - 1
		newWrappedOffset := 0

		for i := len(lines) - 1; i >= 0; i-- {
			wrapped := ta.wrapLine(lines[i])
			if targetDisplayLines+len(wrapped) >= ta.height {
				newViewport = i

				excess := (targetDisplayLines + len(wrapped)) - ta.height
				newWrappedOffset = excess
				break
			}
			targetDisplayLines += len(wrapped)
			newViewport = i
			newWrappedOffset = 0
		}

		ta.viewportLine = newViewport
		ta.viewportWrappedOffset = newWrappedOffset
	}

	if ta.viewportLine < 0 {
		ta.viewportLine = 0
	}

	if ta.viewportLine >= totalLines {
		ta.viewportLine = totalLines - 1
		if ta.viewportLine < 0 {
			ta.viewportLine = 0
		}
	}
}

func (ta *TextArea) scrollToEnd() {
	lines := ta.getLines()
	if len(lines) == 0 {
		ta.viewportLine = 0
		return
	}

	totalDisplayLines := 0
	for _, line := range lines {
		wrapped := ta.wrapLine(line)
		totalDisplayLines += len(wrapped)
	}

	if totalDisplayLines > ta.height {
		targetDisplayLines := 0
		for i := len(lines) - 1; i >= 0; i-- {
			wrapped := ta.wrapLine(lines[i])
			targetDisplayLines += len(wrapped)
			if targetDisplayLines >= ta.height {
				ta.viewportLine = i
				return
			}
		}
	}
	ta.viewportLine = 0
}
