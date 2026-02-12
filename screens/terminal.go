package screens

import (
	"fmt"
	"yoru/models"
	"yoru/screens/popups"
	"yoru/shared"
	"yoru/ssh"
	"yoru/terminal"
	"yoru/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// NewTerminalScreen creates a new terminal screen for a host
func NewTerminalScreen(host *models.Host) *terminalScreen {
	// Calculate terminal dimensions (screen - tab bar)
	width := shared.GlobalState.ScreenWidth
	height := shared.GlobalState.ScreenHeight - 4 // Subtract tab bar height

	return &terminalScreen{
		hostID:          host.ID,
		host:            host,
		emulator:        terminal.NewEmulator(width, height),
		connectionPopup: popups.NewConnectionPopup(),
		connecting:      true,
		keyCaptureMode:  types.KeyCaptureNormal,
	}
}

func (screen *terminalScreen) Init() tea.Cmd {
	// Show connection popup and start SSH connection
	screen.connectionPopup.Show(
		screen.hostID,
		func() {
			screen.connecting = true
			ssh.RetryConnection(screen.hostID)
		},
		func() {
			screen.connectionPopup.Hide()
			screen.shouldClose = true
		},
	)

	return ssh.InitiateConnection(screen.host)
}

func (screen *terminalScreen) Update(msg tea.Msg) (types.Screen, tea.Cmd) {
	// Handle SSH messages first — these arrive while popup is visible
	// and must not be blocked by the popup early return
	switch message := msg.(type) {
	case types.SSHConnectingMsg:
		if message.HostID == screen.hostID {
			screen.connectionPopup.AppendLog(message.Message)
		}
		return screen, nil

	case types.SSHAuthenticatingMsg:
		if message.HostID == screen.hostID {
			screen.connectionPopup.AppendLog(message.Message)
		}
		return screen, nil

	case types.SSHHostKeyMsg:
		if message.HostID == screen.hostID {
			screen.connectionPopup.ShowHostKeyVerification(
				message.Hostname,
				message.Port,
				message.KeyType,
				message.Fingerprint,
				message.ServerKey,
				func() { // onAccept — add to known hosts and continue
					ssh.ContinueAfterHostKeyVerification(screen.hostID, true)
				},
				func() { // onReject — continue without saving
					ssh.ContinueAfterHostKeyVerification(screen.hostID, false)
				},
			)
		}
		return screen, nil

	case types.SSHConnectedMsg:
		if message.HostID == screen.hostID {
			screen.connecting = false
			screen.connected = true
			if connLog, ok := message.ConnectionLog.(*models.ConnectionLog); ok {
				screen.connectionLog = connLog
			}
			screen.connectionPopup.Hide()
			screen.keyCaptureMode = types.KeyCaptureTerminal
			// Resize: full width, height minus tab bar
			screen.emulator.Resize(shared.GlobalState.ScreenWidth, shared.GlobalState.ScreenHeight-1)
			ssh.ResizeTerminal(screen.hostID, shared.GlobalState.ScreenWidth, shared.GlobalState.ScreenHeight-1)
		}
		return screen, nil

	case types.SSHOutputMsg:
		if message.HostID == screen.hostID && screen.connected {
			screen.emulator.Write(message.Data)
		}
		return screen, nil

	case types.SSHErrorMsg:
		if message.HostID == screen.hostID {
			screen.connecting = false
			screen.connectionPopup.ShowError(message.Error)
		}
		return screen, nil

	case types.SSHDisconnectedMsg:
		if message.HostID == screen.hostID {
			screen.connected = false
			ssh.CloseConnection(screen.hostID)
			return screen, func() tea.Msg { return types.CloseTabMsg{} }
		}
		return screen, nil

	case tea.WindowSizeMsg:
		width := message.Width
		height := message.Height - 1 // tab bar
		screen.emulator.Resize(width, height)
		if screen.connected {
			ssh.ResizeTerminal(screen.hostID, width, height)
		}
		return screen, nil

	case tea.MouseMsg:
		if !screen.connectionPopup.IsVisible() && screen.connected {
			switch message.Type {
			case tea.MouseWheelUp:
				screen.emulator.WheelUp()
				return screen, nil
			case tea.MouseWheelDown:
				if screen.emulator.IsScrolled() {
					screen.emulator.WheelDown()
				}
				return screen, nil
			}
		}
		return screen, nil

	case tea.KeyMsg:
		if screen.connectionPopup.IsVisible() {
			screen.connectionPopup.Update(msg)
			if screen.shouldClose {
				screen.shouldClose = false
				ssh.CloseConnection(screen.hostID)
				return screen, func() tea.Msg { return types.CloseTabMsg{} }
			}
			return screen, nil
		}

		if message.Type == tea.KeyCtrlCloseBracket {
			if screen.keyCaptureMode == types.KeyCaptureTerminal {
				screen.keyCaptureMode = types.KeyCaptureNormal
			} else if screen.connected {
				screen.keyCaptureMode = types.KeyCaptureTerminal
			}
			return screen, nil
		}

		if screen.keyCaptureMode == types.KeyCaptureTerminal && screen.connected {
			data := keyToBytes(message)
			if len(data) > 0 {
				ssh.SendInput(screen.hostID, data)
			}
			return screen, nil
		}

		// In normal mode, handle special keys
		if cmd := screen.OnKeyPress(message); cmd != nil {
			return screen, cmd
		}
	}

	return screen, nil
}

func (screen *terminalScreen) View() string {
	// Show connection popup if visible
	if screen.connectionPopup.IsVisible() {
		return screen.connectionPopup.Render()
	}

	// Show terminal if connected
	if screen.connected {
		return screen.emulator.Render()
	}

	// Show connecting message
	width := shared.GlobalState.ScreenWidth
	height := shared.GlobalState.ScreenHeight - 4
	message := fmt.Sprintf("Connecting to %s...", screen.host.Name)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, message)
}

func (screen *terminalScreen) OnKeyPress(key tea.KeyMsg) tea.Cmd {
	// Note: Terminal automatically enters capture mode when connected
	// Shift+Esc releases capture mode (handled in manager)
	return nil
}

// GetKeyCaptureMode returns the current key capture mode
func (screen *terminalScreen) GetKeyCaptureMode() types.KeyCaptureMode {
	return screen.keyCaptureMode
}

// keyToBytes converts a key message to bytes for SSH input
func keyToBytes(key tea.KeyMsg) []byte {
	switch key.Type {
	case tea.KeyEnter:
		return []byte{'\r'}
	case tea.KeyBackspace:
		return []byte{0x7f}
	case tea.KeyTab:
		return []byte{'\t'}
	case tea.KeyEscape:
		return []byte{0x1b}
	case tea.KeyUp:
		return []byte{0x1b, '[', 'A'}
	case tea.KeyDown:
		return []byte{0x1b, '[', 'B'}
	case tea.KeyRight:
		return []byte{0x1b, '[', 'C'}
	case tea.KeyLeft:
		return []byte{0x1b, '[', 'D'}
	case tea.KeyHome:
		return []byte{0x1b, '[', 'H'}
	case tea.KeyEnd:
		return []byte{0x1b, '[', 'F'}
	case tea.KeyPgUp:
		return []byte{0x1b, '[', '5', '~'}
	case tea.KeyPgDown:
		return []byte{0x1b, '[', '6', '~'}
	case tea.KeyDelete:
		return []byte{0x1b, '[', '3', '~'}
	case tea.KeyInsert:
		return []byte{0x1b, '[', '2', '~'}
	case tea.KeyCtrlA:
		return []byte{0x01}
	case tea.KeyCtrlB:
		return []byte{0x02}
	case tea.KeyCtrlC:
		return []byte{0x03}
	case tea.KeyCtrlD:
		return []byte{0x04}
	case tea.KeyCtrlE:
		return []byte{0x05}
	case tea.KeyCtrlF:
		return []byte{0x06}
	case tea.KeyCtrlG:
		return []byte{0x07}
	case tea.KeyCtrlH:
		return []byte{0x08}
	case tea.KeyCtrlJ:
		return []byte{0x0a}
	case tea.KeyCtrlK:
		return []byte{0x0b}
	case tea.KeyCtrlL:
		return []byte{0x0c}
	case tea.KeyCtrlN:
		return []byte{0x0e}
	case tea.KeyCtrlO:
		return []byte{0x0f}
	case tea.KeyCtrlP:
		return []byte{0x10}
	case tea.KeyCtrlR:
		return []byte{0x12}
	case tea.KeyCtrlS:
		return []byte{0x13}
	case tea.KeyCtrlT:
		return []byte{0x14}
	case tea.KeyCtrlU:
		return []byte{0x15}
	case tea.KeyCtrlV:
		return []byte{0x16}
	case tea.KeyCtrlW:
		return []byte{0x17}
	case tea.KeyCtrlX:
		return []byte{0x18}
	case tea.KeyCtrlY:
		return []byte{0x19}
	case tea.KeyCtrlZ:
		return []byte{0x1a}
	case tea.KeyCtrlBackslash:
		return []byte{0x1c}
	case tea.KeyCtrlUnderscore:
		return []byte{0x1f}
	case tea.KeyRunes:
		return []byte(string(key.Runes))
	default:
		str := key.String()
		if len(str) >= 5 && str[:4] == "ctrl" {
			ch := str[5:]
			if len(ch) == 1 {
				c := ch[0]
				if c >= 'a' && c <= 'z' {
					return []byte{byte(c - 'a' + 1)}
				}
			}
		}
		if len(str) > 0 {
			return []byte(str)
		}
	}
	return nil
}