package terminal

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Emulator is a VT100/xterm terminal emulator
type Emulator struct {
	Buffer *Buffer
	Parser *Parser
	Title  string

	// Tab stops (every 8 columns by default)
	TabStops map[int]bool
}

// NewEmulator creates a new terminal emulator
func NewEmulator(width, height int) *Emulator {
	e := &Emulator{
		Buffer:   NewBuffer(width, height),
		Parser:   NewParser(),
		TabStops: make(map[int]bool),
	}

	// Set default tab stops (every 8 columns)
	for i := 0; i < width; i += 8 {
		e.TabStops[i] = true
	}

	return e
}

// Resize resizes the terminal
func (e *Emulator) Resize(width, height int) {
	e.Buffer.Resize(width, height)

	// Reset tab stops for new width
	e.TabStops = make(map[int]bool)
	for i := 0; i < width; i += 8 {
		e.TabStops[i] = true
	}
}

// Write processes input data through the parser and applies actions to the buffer
func (e *Emulator) Write(data []byte) {
	actions := e.Parser.Parse(data)

	for _, action := range actions {
		e.processAction(action)
	}
}

// processAction applies a single action to the buffer
func (e *Emulator) processAction(action Action) {
	switch a := action.(type) {
	case TextAction:
		// Write each character
		for _, r := range a.Text {
			e.Buffer.PutRune(r)
		}

	case BellAction:
		// Visual bell (not implemented)

	case BackspaceAction:
		if e.Buffer.CursorX > 0 {
			e.Buffer.CursorX--
		}

	case TabAction:
		// Move to next tab stop
		startX := e.Buffer.CursorX + 1
		for x := startX; x < e.Buffer.Width; x++ {
			if e.TabStops[x] {
				e.Buffer.CursorX = x
				return
			}
		}
		// No tab stop found, move to end
		e.Buffer.CursorX = e.Buffer.Width - 1

	case LineFeedAction:
		e.Buffer.CursorY++
		if e.Buffer.CursorY > e.Buffer.ScrollBottom {
			e.Buffer.ScrollUp(1)
			e.Buffer.CursorY = e.Buffer.ScrollBottom
		}

	case CarriageReturnAction:
		e.Buffer.CursorX = 0

	case CursorUpAction:
		e.Buffer.MoveCursorRelative(-a.N, 0)

	case CursorDownAction:
		e.Buffer.MoveCursorRelative(a.N, 0)

	case CursorForwardAction:
		e.Buffer.MoveCursorRelative(0, a.N)

	case CursorBackwardAction:
		e.Buffer.MoveCursorRelative(0, -a.N)

	case CursorPositionAction:
		e.Buffer.MoveCursor(a.Row, a.Col)

	case CursorNextLineAction:
		e.Buffer.CursorY += a.N
		e.Buffer.CursorX = 0
		if e.Buffer.CursorY >= e.Buffer.Height {
			e.Buffer.CursorY = e.Buffer.Height - 1
		}

	case CursorPrevLineAction:
		e.Buffer.CursorY -= a.N
		e.Buffer.CursorX = 0
		if e.Buffer.CursorY < 0 {
			e.Buffer.CursorY = 0
		}

	case CursorColumnAction:
		e.Buffer.CursorX = a.Col - 1
		if e.Buffer.CursorX < 0 {
			e.Buffer.CursorX = 0
		}
		if e.Buffer.CursorX >= e.Buffer.Width {
			e.Buffer.CursorX = e.Buffer.Width - 1
		}

	case EraseDisplayAction:
		e.Buffer.EraseDisplay(a.Mode)

	case EraseLineAction:
		e.Buffer.EraseLine(a.Mode)

	case ScrollUpAction:
		e.Buffer.ScrollUp(a.N)

	case ScrollDownAction:
		e.Buffer.ScrollDown(a.N)

	case SGRAction:
		e.Buffer.SetSGR(a.Params)

	case SaveCursorAction:
		e.Buffer.SaveCursor()

	case RestoreCursorAction:
		e.Buffer.RestoreCursor()

	case SetModeAction:
		e.setMode(a.Modes)

	case ResetModeAction:
		e.resetMode(a.Modes)

	case InsertLinesAction:
		e.Buffer.InsertLines(a.N)

	case DeleteLinesAction:
		e.Buffer.DeleteLines(a.N)

	case InsertCharsAction:
		e.Buffer.InsertChars(a.N)

	case DeleteCharsAction:
		e.Buffer.DeleteChars(a.N)

	case SetTitleAction:
		e.Title = a.Title
	}
}

// setMode sets terminal modes
func (e *Emulator) setMode(modes []int) {
	for _, mode := range modes {
		switch mode {
		case 25: // Show cursor
			e.Buffer.CursorVisible = true
		case 7: // Auto-wrap mode
			e.Buffer.AutoWrap = true
		case 6: // Origin mode
			e.Buffer.OriginMode = true
		}
	}
}

// resetMode resets terminal modes
func (e *Emulator) resetMode(modes []int) {
	for _, mode := range modes {
		switch mode {
		case 25: // Hide cursor
			e.Buffer.CursorVisible = false
		case 7: // Auto-wrap mode
			e.Buffer.AutoWrap = false
		case 6: // Origin mode
			e.Buffer.OriginMode = false
		}
	}
}

// Render renders the terminal buffer to a string with ANSI styling
func (e *Emulator) Render() string {
	var result strings.Builder

	for y := 0; y < e.Buffer.Height; y++ {
		for x := 0; x < e.Buffer.Width; x++ {
			cell := e.Buffer.Lines[y][x]
			isCursor := e.Buffer.CursorVisible && x == e.Buffer.CursorX && y == e.Buffer.CursorY

			// Determine the character to render
			ch := " "
			if cell.Rune != 0 && cell.Rune != ' ' && !cell.Hidden {
				ch = string(cell.Rune)
			}

			// Check if this cell needs any styling at all
			hasStyle := isCursor || cell.Foreground >= 0 || cell.Background >= 0 ||
				cell.Bold || cell.Italic || cell.Underline || cell.Reverse ||
				cell.Dim || cell.Strike || cell.Blink

			if !hasStyle {
				result.WriteString(ch)
				continue
			}

			// Build style only for cells that need it
			style := lipgloss.NewStyle()

			if isCursor {
				// Cursor: reverse video block; swap fg/bg with defaults
				fg := cell.Background
				bg := cell.Foreground
				if fg < 0 {
					fg = 7 // default white
				}
				if bg < 0 {
					bg = 0 // default black
				}
				style = style.
					Foreground(lipgloss.Color(colorToString(fg))).
					Background(lipgloss.Color(colorToString(bg)))
			} else {
				// Apply colors
				if cell.Foreground >= 0 {
					style = style.Foreground(lipgloss.Color(colorToString(cell.Foreground)))
				}
				if cell.Background >= 0 {
					style = style.Background(lipgloss.Color(colorToString(cell.Background)))
				}
				// Apply attributes
				if cell.Bold {
					style = style.Bold(true)
				}
				if cell.Italic {
					style = style.Italic(true)
				}
				if cell.Underline {
					style = style.Underline(true)
				}
				if cell.Reverse {
					fg := cell.Background
					bg := cell.Foreground
					if fg >= 0 {
						style = style.Foreground(lipgloss.Color(colorToString(fg)))
					}
					if bg >= 0 {
						style = style.Background(lipgloss.Color(colorToString(bg)))
					}
				}
				if cell.Dim {
					style = style.Faint(true)
				}
				if cell.Strike {
					style = style.Strikethrough(true)
				}
				if cell.Blink {
					style = style.Blink(true)
				}
			}

			result.WriteString(style.Render(ch))
		}
		if y < e.Buffer.Height-1 {
			result.WriteRune('\n')
		}
	}

	return result.String()
}

// colorToString converts ANSI color code to color string
func colorToString(code int) string {
	switch code {
	case 0:
		return "0" // Black
	case 1:
		return "1" // Red
	case 2:
		return "2" // Green
	case 3:
		return "3" // Yellow
	case 4:
		return "4" // Blue
	case 5:
		return "5" // Magenta
	case 6:
		return "6" // Cyan
	case 7:
		return "7" // White
	case 8:
		return "8" // Bright Black
	case 9:
		return "9" // Bright Red
	case 10:
		return "10" // Bright Green
	case 11:
		return "11" // Bright Yellow
	case 12:
		return "12" // Bright Blue
	case 13:
		return "13" // Bright Magenta
	case 14:
		return "14" // Bright Cyan
	case 15:
		return "15" // Bright White
	default:
		// For 256-color codes, use the code directly
		return strconv.Itoa(code)
	}
}

// GetCursorPosition returns the current cursor position
func (e *Emulator) GetCursorPosition() (x, y int) {
	return e.Buffer.CursorX, e.Buffer.CursorY
}

// IsCursorVisible returns whether the cursor is visible
func (e *Emulator) IsCursorVisible() bool {
	return e.Buffer.CursorVisible
}

// Clear clears the terminal
func (e *Emulator) Clear() {
	e.Buffer.Clear()
}