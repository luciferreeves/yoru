package terminal

import (
	"strconv"
	"strings"
)

// Parser state machine states
type parserState int

const (
	stateGround parserState = iota
	stateEscape
	stateCSI
	stateOSC
	stateCSIParam
	stateCSIIntermediate
)

// Action represents a parsed terminal action
type Action interface {
	isAction()
}

// Text action - regular text to display
type TextAction struct {
	Text string
}

func (TextAction) isAction() {}

// Control character actions
type BellAction struct{}
type BackspaceAction struct{}
type TabAction struct{}
type LineFeedAction struct{}
type CarriageReturnAction struct{}

func (BellAction) isAction()          {}
func (BackspaceAction) isAction()     {}
func (TabAction) isAction()           {}
func (LineFeedAction) isAction()      {}
func (CarriageReturnAction) isAction() {}

// Cursor movement
type CursorUpAction struct{ N int }
type CursorDownAction struct{ N int }
type CursorForwardAction struct{ N int }
type CursorBackwardAction struct{ N int }
type CursorPositionAction struct{ Row, Col int }
type CursorNextLineAction struct{ N int }
type CursorPrevLineAction struct{ N int }
type CursorColumnAction struct{ Col int }

func (CursorUpAction) isAction()       {}
func (CursorDownAction) isAction()     {}
func (CursorForwardAction) isAction()  {}
func (CursorBackwardAction) isAction() {}
func (CursorPositionAction) isAction() {}
func (CursorNextLineAction) isAction() {}
func (CursorPrevLineAction) isAction() {}
func (CursorColumnAction) isAction()   {}

// Erase actions
type EraseDisplayAction struct{ Mode int } // 0=below, 1=above, 2=all, 3=scrollback
type EraseLineAction struct{ Mode int }    // 0=right, 1=left, 2=all

func (EraseDisplayAction) isAction() {}
func (EraseLineAction) isAction()    {}

// Scrolling
type ScrollUpAction struct{ N int }
type ScrollDownAction struct{ N int }

func (ScrollUpAction) isAction()   {}
func (ScrollDownAction) isAction() {}

// SGR (Select Graphic Rendition) - colors and styles
type SGRAction struct {
	Params []int
}

func (SGRAction) isAction() {}

// Save/restore cursor
type SaveCursorAction struct{}
type RestoreCursorAction struct{}

func (SaveCursorAction) isAction()    {}
func (RestoreCursorAction) isAction() {}

// Set/reset mode
type SetModeAction struct{ Modes []int }
type ResetModeAction struct{ Modes []int }

func (SetModeAction) isAction()   {}
func (ResetModeAction) isAction() {}

// Insert/delete
type InsertLinesAction struct{ N int }
type DeleteLinesAction struct{ N int }
type InsertCharsAction struct{ N int }
type DeleteCharsAction struct{ N int }

func (InsertLinesAction) isAction() {}
func (DeleteLinesAction) isAction() {}
func (InsertCharsAction) isAction() {}
func (DeleteCharsAction) isAction() {}

// OSC sequences
type SetTitleAction struct{ Title string }

func (SetTitleAction) isAction() {}

// Parser parses VT100/xterm escape sequences
type Parser struct {
	state        parserState
	params       []int
	intermediate []byte
	buffer       []byte
}

// NewParser creates a new parser
func NewParser() *Parser {
	return &Parser{
		state:  stateGround,
		params: make([]int, 0, 16),
	}
}

// Parse parses input data and returns a slice of actions
func (p *Parser) Parse(data []byte) []Action {
	actions := make([]Action, 0)
	textStart := -1

	for i := 0; i < len(data); i++ {
		b := data[i]

		switch p.state {
		case stateGround:
			switch b {
			case 0x07: // BEL
				if textStart >= 0 {
					actions = append(actions, TextAction{Text: string(data[textStart:i])})
					textStart = -1
				}
				actions = append(actions, BellAction{})
			case 0x08: // BS
				if textStart >= 0 {
					actions = append(actions, TextAction{Text: string(data[textStart:i])})
					textStart = -1
				}
				actions = append(actions, BackspaceAction{})
			case 0x09: // TAB
				if textStart >= 0 {
					actions = append(actions, TextAction{Text: string(data[textStart:i])})
					textStart = -1
				}
				actions = append(actions, TabAction{})
			case 0x0A, 0x0B, 0x0C: // LF, VT, FF
				if textStart >= 0 {
					actions = append(actions, TextAction{Text: string(data[textStart:i])})
					textStart = -1
				}
				actions = append(actions, LineFeedAction{})
			case 0x0D: // CR
				if textStart >= 0 {
					actions = append(actions, TextAction{Text: string(data[textStart:i])})
					textStart = -1
				}
				actions = append(actions, CarriageReturnAction{})
			case 0x1B: // ESC
				if textStart >= 0 {
					actions = append(actions, TextAction{Text: string(data[textStart:i])})
					textStart = -1
				}
				p.state = stateEscape
				p.params = p.params[:0]
				p.intermediate = p.intermediate[:0]
				p.buffer = p.buffer[:0]
			default:
				// Regular text
				if textStart < 0 {
					textStart = i
				}
			}

		case stateEscape:
			switch b {
			case '[': // CSI
				p.state = stateCSI
			case ']': // OSC
				p.state = stateOSC
			case '7': // Save cursor (DECSC)
				actions = append(actions, SaveCursorAction{})
				p.state = stateGround
			case '8': // Restore cursor (DECRC)
				actions = append(actions, RestoreCursorAction{})
				p.state = stateGround
			case 'M': // Reverse index (scroll down)
				actions = append(actions, ScrollDownAction{N: 1})
				p.state = stateGround
			case 'D': // Index (scroll up)
				actions = append(actions, ScrollUpAction{N: 1})
				p.state = stateGround
			default:
				// Unknown escape sequence, ignore
				p.state = stateGround
			}

		case stateCSI:
			if b >= '0' && b <= '9' {
				p.state = stateCSIParam
				p.buffer = append(p.buffer, b)
			} else if b == ';' {
				p.params = append(p.params, 0)
				p.state = stateCSIParam
			} else if b == '?' || b == '>' || b == '<' {
				// Private-mode prefix (e.g. ESC[?25h) — skip, continue to params
				p.state = stateCSIParam
			} else {
				// No parameters, process command
				action := p.processCSI(b)
				if action != nil {
					actions = append(actions, action)
				}
				p.state = stateGround
			}

		case stateCSIParam:
			if b >= '0' && b <= '9' {
				p.buffer = append(p.buffer, b)
			} else if b == ';' {
				// Parse accumulated number
				if len(p.buffer) > 0 {
					n, _ := strconv.Atoi(string(p.buffer))
					p.params = append(p.params, n)
					p.buffer = p.buffer[:0]
				} else {
					p.params = append(p.params, 0)
				}
			} else if b >= 0x20 && b <= 0x2F {
				// Intermediate byte
				if len(p.buffer) > 0 {
					n, _ := strconv.Atoi(string(p.buffer))
					p.params = append(p.params, n)
					p.buffer = p.buffer[:0]
				}
				p.intermediate = append(p.intermediate, b)
				p.state = stateCSIIntermediate
			} else {
				// Final byte - process command
				if len(p.buffer) > 0 {
					n, _ := strconv.Atoi(string(p.buffer))
					p.params = append(p.params, n)
					p.buffer = p.buffer[:0]
				}
				action := p.processCSI(b)
				if action != nil {
					actions = append(actions, action)
				}
				p.state = stateGround
			}

		case stateCSIIntermediate:
			if b >= 0x20 && b <= 0x2F {
				p.intermediate = append(p.intermediate, b)
			} else {
				// Final byte
				action := p.processCSI(b)
				if action != nil {
					actions = append(actions, action)
				}
				p.state = stateGround
			}

		case stateOSC:
			if b == 0x07 || (b == 0x1B && i+1 < len(data) && data[i+1] == '\\') {
				// OSC terminator (BEL or ESC \)
				action := p.processOSC()
				if action != nil {
					actions = append(actions, action)
				}
				if b == 0x1B {
					i++ // Skip the backslash
				}
				p.state = stateGround
			} else {
				p.buffer = append(p.buffer, b)
			}
		}
	}

	// Flush any remaining text
	if textStart >= 0 {
		actions = append(actions, TextAction{Text: string(data[textStart:])})
	}

	return actions
}

// processCSI processes a CSI sequence
func (p *Parser) processCSI(final byte) Action {
	// Get first parameter with default
	n := 1
	if len(p.params) > 0 {
		n = p.params[0]
		if n == 0 {
			n = 1
		}
	}

	switch final {
	case 'A': // Cursor Up
		return CursorUpAction{N: n}
	case 'B': // Cursor Down
		return CursorDownAction{N: n}
	case 'C': // Cursor Forward
		return CursorForwardAction{N: n}
	case 'D': // Cursor Backward
		return CursorBackwardAction{N: n}
	case 'E': // Cursor Next Line
		return CursorNextLineAction{N: n}
	case 'F': // Cursor Previous Line
		return CursorPrevLineAction{N: n}
	case 'G': // Cursor Horizontal Absolute
		return CursorColumnAction{Col: n}
	case 'H', 'f': // Cursor Position
		row := 1
		col := 1
		if len(p.params) >= 1 {
			row = p.params[0]
			if row == 0 {
				row = 1
			}
		}
		if len(p.params) >= 2 {
			col = p.params[1]
			if col == 0 {
				col = 1
			}
		}
		return CursorPositionAction{Row: row, Col: col}
	case 'J': // Erase Display
		mode := 0
		if len(p.params) > 0 {
			mode = p.params[0]
		}
		return EraseDisplayAction{Mode: mode}
	case 'K': // Erase Line
		mode := 0
		if len(p.params) > 0 {
			mode = p.params[0]
		}
		return EraseLineAction{Mode: mode}
	case 'L': // Insert Lines
		return InsertLinesAction{N: n}
	case 'M': // Delete Lines
		return DeleteLinesAction{N: n}
	case 'P': // Delete Characters
		return DeleteCharsAction{N: n}
	case 'S': // Scroll Up
		return ScrollUpAction{N: n}
	case 'T': // Scroll Down
		return ScrollDownAction{N: n}
	case '@': // Insert Characters
		return InsertCharsAction{N: n}
	case 'm': // SGR
		if len(p.params) == 0 {
			p.params = append(p.params, 0)
		}
		return SGRAction{Params: p.params}
	case 'h': // Set Mode
		return SetModeAction{Modes: p.params}
	case 'l': // Reset Mode
		return ResetModeAction{Modes: p.params}
	case 's': // Save cursor position (ANSI)
		return SaveCursorAction{}
	case 'u': // Restore cursor position (ANSI)
		return RestoreCursorAction{}
	default:
		// Unknown CSI sequence
		return nil
	}
}

// processOSC processes an OSC sequence
func (p *Parser) processOSC() Action {
	if len(p.buffer) == 0 {
		return nil
	}

	// Parse OSC command: Ps ; Pt
	parts := strings.SplitN(string(p.buffer), ";", 2)
	if len(parts) < 2 {
		return nil
	}

	ps, _ := strconv.Atoi(parts[0])
	pt := parts[1]

	switch ps {
	case 0, 2: // Set window title
		return SetTitleAction{Title: pt}
	default:
		return nil
	}
}