package terminal

import (
	"strings"
)

// Cell represents a single terminal cell
type Cell struct {
	Rune       rune
	Foreground int  // ANSI color code
	Background int  // ANSI color code
	Bold       bool
	Dim        bool
	Italic     bool
	Underline  bool
	Blink      bool
	Reverse    bool
	Hidden     bool
	Strike     bool
}

// Buffer represents the terminal screen buffer
type Buffer struct {
	Width  int
	Height int

	// Screen lines (visible area)
	Lines [][]Cell

	// Scrollback buffer
	Scrollback [][]Cell
	MaxScrollback int

	// Cursor position (0-indexed)
	CursorX int
	CursorY int

	// Saved cursor position
	SavedCursorX int
	SavedCursorY int
	SavedAttrs   Cell

	// Current cell attributes
	CurrentAttrs Cell

	// Scroll region
	ScrollTop    int
	ScrollBottom int

	// Modes
	CursorVisible bool
	AutoWrap      bool
	OriginMode    bool
}

// NewBuffer creates a new terminal buffer
func NewBuffer(width, height int) *Buffer {
	b := &Buffer{
		Width:         width,
		Height:        height,
		Lines:         make([][]Cell, height),
		Scrollback:    make([][]Cell, 0),
		MaxScrollback: 1000,
		CursorVisible: true,
		AutoWrap:      true,
		ScrollTop:     0,
		ScrollBottom:  height - 1,
	}

	blank := Cell{
		Rune:       ' ',
		Foreground: -1,
		Background: -1,
	}

	for i := 0; i < height; i++ {
		b.Lines[i] = make([]Cell, width)
		for j := 0; j < width; j++ {
			b.Lines[i][j] = blank
		}
	}

	// Default attributes
	b.CurrentAttrs = Cell{
		Rune:       ' ',
		Foreground: -1, // Default color
		Background: -1,
	}

	return b
}

// Resize resizes the buffer
func (b *Buffer) Resize(width, height int) {
	if width == b.Width && height == b.Height {
		return
	}

	blank := Cell{
		Rune:       ' ',
		Foreground: -1,
		Background: -1,
	}

	newLines := make([][]Cell, height)
	for i := 0; i < height; i++ {
		newLines[i] = make([]Cell, width)
		for j := 0; j < width; j++ {
			newLines[i][j] = blank
		}
		if i < len(b.Lines) {
			copy(newLines[i], b.Lines[i])
		}
	}

	b.Lines = newLines
	b.Width = width
	b.Height = height
	b.ScrollBottom = height - 1

	// Clamp cursor
	if b.CursorX >= width {
		b.CursorX = width - 1
	}
	if b.CursorY >= height {
		b.CursorY = height - 1
	}
}

// PutRune writes a rune at the current cursor position
func (b *Buffer) PutRune(r rune) {
	if b.CursorY < 0 || b.CursorY >= b.Height {
		return
	}
	if b.CursorX < 0 || b.CursorX >= b.Width {
		return
	}

	cell := b.CurrentAttrs
	cell.Rune = r
	b.Lines[b.CursorY][b.CursorX] = cell

	// Advance cursor
	b.CursorX++
	if b.CursorX >= b.Width {
		if b.AutoWrap {
			b.CursorX = 0
			b.CursorY++
			if b.CursorY > b.ScrollBottom {
				b.ScrollUp(1)
				b.CursorY = b.ScrollBottom
			}
		} else {
			b.CursorX = b.Width - 1
		}
	}
}

// MoveCursor moves the cursor to an absolute position (1-indexed)
func (b *Buffer) MoveCursor(row, col int) {
	// Convert to 0-indexed
	row--
	col--

	// Handle origin mode
	if b.OriginMode {
		row += b.ScrollTop
	}

	// Clamp to valid range
	if row < 0 {
		row = 0
	}
	if row >= b.Height {
		row = b.Height - 1
	}
	if col < 0 {
		col = 0
	}
	if col >= b.Width {
		col = b.Width - 1
	}

	b.CursorY = row
	b.CursorX = col
}

// MoveCursorRelative moves the cursor relative to current position
func (b *Buffer) MoveCursorRelative(dy, dx int) {
	b.CursorY += dy
	b.CursorX += dx

	// Clamp
	if b.CursorY < 0 {
		b.CursorY = 0
	}
	if b.CursorY >= b.Height {
		b.CursorY = b.Height - 1
	}
	if b.CursorX < 0 {
		b.CursorX = 0
	}
	if b.CursorX >= b.Width {
		b.CursorX = b.Width - 1
	}
}

// EraseDisplay erases portions of the display
func (b *Buffer) EraseDisplay(mode int) {
	blank := b.CurrentAttrs
	blank.Rune = ' '

	switch mode {
	case 0: // Erase from cursor to end of screen
		// Clear rest of current line
		for x := b.CursorX; x < b.Width; x++ {
			b.Lines[b.CursorY][x] = blank
		}
		// Clear lines below
		for y := b.CursorY + 1; y < b.Height; y++ {
			for x := 0; x < b.Width; x++ {
				b.Lines[y][x] = blank
			}
		}

	case 1: // Erase from beginning of screen to cursor
		// Clear lines above
		for y := 0; y < b.CursorY; y++ {
			for x := 0; x < b.Width; x++ {
				b.Lines[y][x] = blank
			}
		}
		// Clear beginning of current line
		for x := 0; x <= b.CursorX; x++ {
			b.Lines[b.CursorY][x] = blank
		}

	case 2, 3: // Erase entire screen (3 also clears scrollback)
		for y := 0; y < b.Height; y++ {
			for x := 0; x < b.Width; x++ {
				b.Lines[y][x] = blank
			}
		}
		if mode == 3 {
			b.Scrollback = b.Scrollback[:0]
		}
	}
}

// EraseLine erases portions of the current line
func (b *Buffer) EraseLine(mode int) {
	if b.CursorY < 0 || b.CursorY >= b.Height {
		return
	}

	blank := b.CurrentAttrs
	blank.Rune = ' '

	switch mode {
	case 0: // Erase from cursor to end of line
		for x := b.CursorX; x < b.Width; x++ {
			b.Lines[b.CursorY][x] = blank
		}

	case 1: // Erase from beginning of line to cursor
		for x := 0; x <= b.CursorX; x++ {
			b.Lines[b.CursorY][x] = blank
		}

	case 2: // Erase entire line
		for x := 0; x < b.Width; x++ {
			b.Lines[b.CursorY][x] = blank
		}
	}
}

// ScrollUp scrolls the scroll region up by n lines
func (b *Buffer) ScrollUp(n int) {
	if n <= 0 {
		return
	}

	// Move lines to scrollback
	for i := 0; i < n && b.ScrollTop < len(b.Lines); i++ {
		if b.ScrollTop < len(b.Lines) {
			line := b.Lines[b.ScrollTop]
			b.Scrollback = append(b.Scrollback, line)
			if len(b.Scrollback) > b.MaxScrollback {
				b.Scrollback = b.Scrollback[1:]
			}
		}
	}

	// Shift lines up
	for y := b.ScrollTop; y <= b.ScrollBottom-n; y++ {
		if y+n <= b.ScrollBottom {
			copy(b.Lines[y], b.Lines[y+n])
		}
	}

	// Clear bottom lines
	blank := b.CurrentAttrs
	blank.Rune = ' '
	for y := b.ScrollBottom - n + 1; y <= b.ScrollBottom; y++ {
		if y >= 0 && y < b.Height {
			for x := 0; x < b.Width; x++ {
				b.Lines[y][x] = blank
			}
		}
	}
}

// ScrollDown scrolls the scroll region down by n lines
func (b *Buffer) ScrollDown(n int) {
	if n <= 0 {
		return
	}

	// Shift lines down
	for y := b.ScrollBottom; y >= b.ScrollTop+n; y-- {
		if y-n >= b.ScrollTop {
			copy(b.Lines[y], b.Lines[y-n])
		}
	}

	// Clear top lines
	blank := b.CurrentAttrs
	blank.Rune = ' '
	for y := b.ScrollTop; y < b.ScrollTop+n && y <= b.ScrollBottom; y++ {
		if y >= 0 && y < b.Height {
			for x := 0; x < b.Width; x++ {
				b.Lines[y][x] = blank
			}
		}
	}
}

// InsertLines inserts n blank lines at cursor position
func (b *Buffer) InsertLines(n int) {
	if b.CursorY < b.ScrollTop || b.CursorY > b.ScrollBottom {
		return
	}

	// Shift lines down
	for y := b.ScrollBottom; y >= b.CursorY+n; y-- {
		if y-n >= b.CursorY {
			copy(b.Lines[y], b.Lines[y-n])
		}
	}

	// Clear inserted lines
	blank := b.CurrentAttrs
	blank.Rune = ' '
	for y := b.CursorY; y < b.CursorY+n && y <= b.ScrollBottom; y++ {
		for x := 0; x < b.Width; x++ {
			b.Lines[y][x] = blank
		}
	}
}

// DeleteLines deletes n lines at cursor position
func (b *Buffer) DeleteLines(n int) {
	if b.CursorY < b.ScrollTop || b.CursorY > b.ScrollBottom {
		return
	}

	// Shift lines up
	for y := b.CursorY; y <= b.ScrollBottom-n; y++ {
		if y+n <= b.ScrollBottom {
			copy(b.Lines[y], b.Lines[y+n])
		}
	}

	// Clear bottom lines
	blank := b.CurrentAttrs
	blank.Rune = ' '
	for y := b.ScrollBottom - n + 1; y <= b.ScrollBottom; y++ {
		if y >= 0 && y < b.Height {
			for x := 0; x < b.Width; x++ {
				b.Lines[y][x] = blank
			}
		}
	}
}

// InsertChars inserts n blank characters at cursor position
func (b *Buffer) InsertChars(n int) {
	if b.CursorY < 0 || b.CursorY >= b.Height {
		return
	}

	// Shift characters right
	line := b.Lines[b.CursorY]
	for x := b.Width - 1; x >= b.CursorX+n; x-- {
		if x-n >= b.CursorX {
			line[x] = line[x-n]
		}
	}

	// Clear inserted characters
	blank := b.CurrentAttrs
	blank.Rune = ' '
	for x := b.CursorX; x < b.CursorX+n && x < b.Width; x++ {
		line[x] = blank
	}
}

// DeleteChars deletes n characters at cursor position
func (b *Buffer) DeleteChars(n int) {
	if b.CursorY < 0 || b.CursorY >= b.Height {
		return
	}

	// Shift characters left
	line := b.Lines[b.CursorY]
	for x := b.CursorX; x < b.Width-n; x++ {
		if x+n < b.Width {
			line[x] = line[x+n]
		}
	}

	// Clear end characters
	blank := b.CurrentAttrs
	blank.Rune = ' '
	for x := b.Width - n; x < b.Width; x++ {
		if x >= 0 {
			line[x] = blank
		}
	}
}

// SaveCursor saves the current cursor position and attributes
func (b *Buffer) SaveCursor() {
	b.SavedCursorX = b.CursorX
	b.SavedCursorY = b.CursorY
	b.SavedAttrs = b.CurrentAttrs
}

// RestoreCursor restores the saved cursor position and attributes
func (b *Buffer) RestoreCursor() {
	b.CursorX = b.SavedCursorX
	b.CursorY = b.SavedCursorY
	b.CurrentAttrs = b.SavedAttrs
}

// SetSGR applies SGR (Select Graphic Rendition) parameters
func (b *Buffer) SetSGR(params []int) {
	for i := 0; i < len(params); i++ {
		param := params[i]

		switch param {
		case 0: // Reset
			b.CurrentAttrs = Cell{
				Rune:       ' ',
				Foreground: -1,
				Background: -1,
			}

		case 1: // Bold
			b.CurrentAttrs.Bold = true
		case 2: // Dim
			b.CurrentAttrs.Dim = true
		case 3: // Italic
			b.CurrentAttrs.Italic = true
		case 4: // Underline
			b.CurrentAttrs.Underline = true
		case 5: // Blink
			b.CurrentAttrs.Blink = true
		case 7: // Reverse
			b.CurrentAttrs.Reverse = true
		case 8: // Hidden
			b.CurrentAttrs.Hidden = true
		case 9: // Strike
			b.CurrentAttrs.Strike = true

		case 22: // Normal intensity
			b.CurrentAttrs.Bold = false
			b.CurrentAttrs.Dim = false
		case 23: // Not italic
			b.CurrentAttrs.Italic = false
		case 24: // Not underline
			b.CurrentAttrs.Underline = false
		case 25: // Not blink
			b.CurrentAttrs.Blink = false
		case 27: // Not reverse
			b.CurrentAttrs.Reverse = false
		case 28: // Not hidden
			b.CurrentAttrs.Hidden = false
		case 29: // Not strike
			b.CurrentAttrs.Strike = false

		case 30, 31, 32, 33, 34, 35, 36, 37: // Foreground colors
			b.CurrentAttrs.Foreground = param - 30
		case 39: // Default foreground
			b.CurrentAttrs.Foreground = -1
		case 40, 41, 42, 43, 44, 45, 46, 47: // Background colors
			b.CurrentAttrs.Background = param - 40
		case 49: // Default background
			b.CurrentAttrs.Background = -1

		case 90, 91, 92, 93, 94, 95, 96, 97: // Bright foreground colors
			b.CurrentAttrs.Foreground = param - 90 + 8
		case 100, 101, 102, 103, 104, 105, 106, 107: // Bright background colors
			b.CurrentAttrs.Background = param - 100 + 8

		case 38: // 256-color or RGB foreground
			if i+2 < len(params) && params[i+1] == 5 {
				// 256-color: ESC[38;5;Nm
				b.CurrentAttrs.Foreground = params[i+2]
				i += 2
			} else if i+4 < len(params) && params[i+1] == 2 {
				// RGB: ESC[38;2;R;G;Bm (not fully supported, map to closest)
				i += 4
			}
		case 48: // 256-color or RGB background
			if i+2 < len(params) && params[i+1] == 5 {
				// 256-color: ESC[48;5;Nm
				b.CurrentAttrs.Background = params[i+2]
				i += 2
			} else if i+4 < len(params) && params[i+1] == 2 {
				// RGB: ESC[48;2;R;G;Bm (not fully supported, map to closest)
				i += 4
			}
		}
	}
}

// GetLine returns a line as a string
func (b *Buffer) GetLine(y int) string {
	if y < 0 || y >= b.Height {
		return ""
	}

	var sb strings.Builder
	for x := 0; x < b.Width; x++ {
		sb.WriteRune(b.Lines[y][x].Rune)
	}
	return sb.String()
}

// Clear clears the entire buffer
func (b *Buffer) Clear() {
	blank := Cell{
		Rune:       ' ',
		Foreground: -1,
		Background: -1,
	}

	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			b.Lines[y][x] = blank
		}
	}

	b.CursorX = 0
	b.CursorY = 0
}