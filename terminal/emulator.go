package terminal

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/hinshun/vt10x"
)

const (
	attrReverse   = 1 << 0
	attrUnderline = 1 << 1
	attrBold      = 1 << 2
	attrGfx       = 1 << 3
	attrItalic    = 1 << 4
	attrBlink     = 1 << 5
)

type Emulator struct {
	vt           vt10x.Terminal
	width        int
	height       int
	scrollOffset int
}

func NewEmulator(width, height int) *Emulator {
	vt := vt10x.New(vt10x.WithSize(width, height))
	return &Emulator{
		vt:           vt,
		width:        width,
		height:       height,
		scrollOffset: 0,
	}
}

func (e *Emulator) Resize(width, height int) {
	e.width = width
	e.height = height
	e.vt.Resize(width, height)
}

func (e *Emulator) Write(data []byte) {
	e.vt.Write(data)
	e.scrollOffset = 0
}

func (e *Emulator) IsScrolled() bool {
	return e.scrollOffset > 0
}

func (e *Emulator) WheelUp() {
	e.scrollOffset += 3
	if e.scrollOffset > 1000 {
		e.scrollOffset = 1000
	}
}

func (e *Emulator) WheelDown() {
	e.scrollOffset -= 3
	if e.scrollOffset < 0 {
		e.scrollOffset = 0
	}
}

func (e *Emulator) Render() string {
	var result strings.Builder
	result.Grow(e.width * e.height * 2)

	cursor := e.vt.Cursor()
	cursorVisible := e.vt.CursorVisible()

	for y := 0; y < e.height; y++ {
		for x := 0; x < e.width; x++ {
			cell := e.vt.Cell(x, y)
			ch := cell.Char

			if ch == 0 {
				ch = ' '
			}

			isCursor := cursorVisible && x == cursor.X && y == cursor.Y

			fg := cell.FG
			bg := cell.BG
			mode := cell.Mode

			if mode&attrReverse != 0 || isCursor {
				fg, bg = bg, fg
			}

			hasStyle := isCursor || fg != vt10x.DefaultFG || bg != vt10x.DefaultBG ||
				mode&(attrBold|attrUnderline|attrItalic|attrBlink) != 0

			if !hasStyle {
				result.WriteRune(ch)
				continue
			}

			style := lipgloss.NewStyle()

			if fg != vt10x.DefaultFG {
				style = style.Foreground(lipgloss.Color(colorToString(fg)))
			}
			if bg != vt10x.DefaultBG {
				style = style.Background(lipgloss.Color(colorToString(bg)))
			}
			if mode&attrBold != 0 {
				style = style.Bold(true)
			}
			if mode&attrItalic != 0 {
				style = style.Italic(true)
			}
			if mode&attrUnderline != 0 {
				style = style.Underline(true)
			}
			if mode&attrBlink != 0 {
				style = style.Blink(true)
			}

			result.WriteString(style.Render(string(ch)))
		}
		if y < e.height-1 {
			result.WriteRune('\n')
		}
	}

	return result.String()
}

func colorToString(color vt10x.Color) string {
	if color >= 0 && color <= 255 {
		return string(rune('0' + int(color)))
	}
	return "7"
}