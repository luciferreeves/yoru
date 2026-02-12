package popups

import (
	"fmt"
	"strings"
	"yoru/repository"
	"yoru/screens/components"
	"yoru/screens/styles"
	"yoru/shared"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	sshlib "golang.org/x/crypto/ssh"
)

type ConnectionState int

const (
	StateConnecting ConnectionState = iota
	StateVerifyingHost
	StateError
	StateConfirmClose
)

type ConnectionPopup struct {
	popup       *components.Popup
	hostID      uint
	logs        []string
	state       ConnectionState
	errorMsg    string
	selectedBtn int

	hostname    string
	port        int
	keyType     string
	fingerprint string
	serverKey   sshlib.PublicKey

	logBoxMinWidth int

	onRetry         func()
	onCancel        func()
	onAcceptHostKey func()
	onRejectHostKey func()
	onConfirmClose  func()
	onCancelClose   func()
}

func NewConnectionPopup() *ConnectionPopup {
	return &ConnectionPopup{
		popup: components.NewPopup(),
		logs:  []string{},
		state: StateConnecting,
	}
}

func (cp *ConnectionPopup) Show(hostID uint, onRetry func(), onCancel func()) {
	cp.hostID = hostID
	cp.logs = []string{}
	cp.state = StateConnecting

	host, _ := repository.GetHostByID(hostID)
	minWidth := 0
	if host != nil {
		verifyWidth := len(fmt.Sprintf("Host '%s:%d' is not in known hosts.", host.Hostname, host.Port))
		if verifyWidth > minWidth {
			minWidth = verifyWidth
		}
	}
	fingerprintWidth := 60
	if fingerprintWidth > minWidth {
		minWidth = fingerprintWidth
	}
	sectionWidth := len("Do you want to add this host to known hosts?")
	if sectionWidth > minWidth {
		minWidth = sectionWidth
	}
	cp.logBoxMinWidth = minWidth

	cp.onRetry = onRetry
	cp.onCancel = onCancel
	cp.popup.SetWidth(shared.GlobalState.ScreenWidth - 20)
	cp.popup.SetBorderless(true)
	cp.popup.SetHeightOffset(1)
	cp.popup.Show(cp.buildContent(), cp.handleInput)
}

func (cp *ConnectionPopup) AppendLog(message string) {
	cp.logs = append(cp.logs, message)
	if len(cp.logs) > 5 {
		cp.logs = cp.logs[len(cp.logs)-5:]
	}
	cp.popup.SetContent(cp.buildContent())
}

func (cp *ConnectionPopup) ShowError(err error) {
	cp.state = StateError
	cp.errorMsg = err.Error()
	cp.selectedBtn = 0
	cp.popup.SetContent(cp.buildContent())
}

func (cp *ConnectionPopup) ShowHostKeyVerification(
	hostname string, port int, keyType string, fingerprint string,
	serverKey sshlib.PublicKey, onAccept func(), onReject func(),
) {
	cp.state = StateVerifyingHost
	cp.hostname = hostname
	cp.port = port
	cp.keyType = keyType
	cp.fingerprint = fingerprint
	cp.serverKey = serverKey
	cp.onAcceptHostKey = onAccept
	cp.onRejectHostKey = onReject
	cp.selectedBtn = 1
	cp.popup.SetContent(cp.buildContent())
}

func (cp *ConnectionPopup) ShowCloseConfirmation(onConfirm func(), onCancelClose func()) {
	cp.state = StateConfirmClose
	cp.onConfirmClose = onConfirm
	cp.onCancelClose = onCancelClose
	cp.selectedBtn = 1
	cp.popup.SetContent(cp.buildContent())
}

func (cp *ConnectionPopup) Hide() {
	cp.popup.Hide()
}

func (cp *ConnectionPopup) IsVisible() bool {
	return cp.popup.IsVisible()
}

func (cp *ConnectionPopup) Update(msg tea.Msg) {
	cp.popup.Update(msg)
}

func (cp *ConnectionPopup) Render() string {
	return cp.popup.Render()
}

func (cp *ConnectionPopup) handleInput(msg tea.Msg) bool {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return false
	}

	switch cp.state {
	case StateError:
		return cp.handleErrorInput(keyMsg)
	case StateVerifyingHost:
		return cp.handleHostKeyInput(keyMsg)
	case StateConfirmClose:
		return cp.handleCloseConfirmInput(keyMsg)
	}
	return false
}

func (cp *ConnectionPopup) handleErrorInput(keyMsg tea.KeyMsg) bool {
	switch keyMsg.String() {
	case "left", "h":
		cp.selectedBtn = 0
		cp.popup.SetContent(cp.buildContent())
	case "right", "l":
		cp.selectedBtn = 1
		cp.popup.SetContent(cp.buildContent())
	case "enter":
		if cp.selectedBtn == 0 && cp.onRetry != nil {
			cp.onRetry()
		} else if cp.onCancel != nil {
			cp.onCancel()
		}
	case "esc":
		if cp.onCancel != nil {
			cp.onCancel()
		}
	default:
		return false
	}
	return true
}

func (cp *ConnectionPopup) handleHostKeyInput(keyMsg tea.KeyMsg) bool {
	switch keyMsg.String() {
	case "left", "h":
		cp.selectedBtn = 0
		cp.popup.SetContent(cp.buildContent())
	case "right", "l":
		cp.selectedBtn = 1
		cp.popup.SetContent(cp.buildContent())
	case "enter":
		if cp.selectedBtn == 0 {
			if cp.onAcceptHostKey != nil {
				cp.onAcceptHostKey()
			}
		} else if cp.onRejectHostKey != nil {
			cp.onRejectHostKey()
		}
		cp.state = StateConnecting
		cp.popup.SetContent(cp.buildContent())
	case "y", "Y":
		if cp.onAcceptHostKey != nil {
			cp.onAcceptHostKey()
		}
		cp.state = StateConnecting
		cp.popup.SetContent(cp.buildContent())
	case "n", "N", "esc":
		if cp.onRejectHostKey != nil {
			cp.onRejectHostKey()
		}
		cp.state = StateConnecting
		cp.popup.SetContent(cp.buildContent())
	default:
		return false
	}
	return true
}

func (cp *ConnectionPopup) handleCloseConfirmInput(keyMsg tea.KeyMsg) bool {
	switch keyMsg.String() {
	case "left", "h":
		cp.selectedBtn = 0
		cp.popup.SetContent(cp.buildContent())
	case "right", "l":
		cp.selectedBtn = 1
		cp.popup.SetContent(cp.buildContent())
	case "enter":
		if cp.selectedBtn == 0 && cp.onConfirmClose != nil {
			cp.onConfirmClose()
		} else if cp.onCancelClose != nil {
			cp.onCancelClose()
		}
	case "y", "Y":
		if cp.onConfirmClose != nil {
			cp.onConfirmClose()
		}
	case "n", "N", "esc":
		if cp.onCancelClose != nil {
			cp.onCancelClose()
		}
	default:
		return false
	}
	return true
}

func (cp *ConnectionPopup) buildContent() string {
	host, _ := repository.GetHostByID(cp.hostID)
	title := "SSH Connection"
	if host != nil {
		title = "SSH Connection: " + host.Name
	}

	parts := []string{styles.PopupTitle.Render(title)}

	// Logs in a bordered box (5 lines viewport)
	logLines := make([]string, 5)
	for i := 0; i < 5; i++ {
		if i < len(cp.logs) {
			logLines[i] = cp.logs[i]
		} else {
			logLines[i] = ""
		}
	}

	maxWidth := cp.logBoxMinWidth

	for _, line := range logLines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	if host != nil {
		verifyWidth := len(fmt.Sprintf("Host '%s:%d' is not in known hosts.", host.Hostname, host.Port))
		if verifyWidth > maxWidth {
			maxWidth = verifyWidth
		}
	}
	if cp.keyType != "" {
		keyWidth := len(fmt.Sprintf("%s key fingerprint is:", cp.keyType))
		if keyWidth > maxWidth {
			maxWidth = keyWidth
		}
		if len(cp.fingerprint) > maxWidth {
			maxWidth = len(cp.fingerprint)
		}
	}
	sectionWidth := len("Do you want to add this host to known hosts?")
	if sectionWidth > maxWidth {
		maxWidth = sectionWidth
	}

	cp.logBoxMinWidth = maxWidth

	for i, line := range logLines {
		if len(line) < maxWidth {
			logLines[i] = line + strings.Repeat(" ", maxWidth-len(line))
		}
	}

	parts = append(parts, styles.PopupLogBox.Render(lipgloss.JoinVertical(lipgloss.Left, logLines...)))

	// Append state-specific content below the log box
	switch cp.state {
	case StateError:
		parts = append(parts,
			styles.PopupError.Render("Error: "+cp.errorMsg),
			styles.PopupButtonsContainer.Render(cp.buildButtons("Retry", "Cancel")))

	case StateVerifyingHost:
		parts = append(parts,
			styles.PopupText.Render(fmt.Sprintf("Host '%s:%d' is not in known hosts.", cp.hostname, cp.port)),
			styles.PopupText.Render(fmt.Sprintf("%s key fingerprint is:", cp.keyType)),
			styles.PopupTextBold.Render(cp.fingerprint),
			styles.PopupSection.Render("Do you want to add this host to known hosts?"),
			styles.PopupButtonsContainer.Render(cp.buildButtons("Yes (y)", "No (n)")))

	case StateConfirmClose:
		parts = append(parts,
			styles.PopupText.Render("An active SSH session is running."),
			styles.PopupSection.Render("Are you sure you want to close this tab?"),
			styles.PopupButtonsContainer.Render(cp.buildButtons("Yes", "No")))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, parts...)
	return styles.PopupContentContainer.Render(content)
}

func (cp *ConnectionPopup) buildButtons(yesLabel, noLabel string) string {
	yesPrefix := "  "
	noPrefix := "  "
	if cp.selectedBtn == 0 {
		yesPrefix = "> "
	} else {
		noPrefix = "> "
	}

	yesButton := styles.PopupButtonYes.Render(yesPrefix + yesLabel)
	noButton := styles.PopupButtonNo.Render(noPrefix + noLabel)

	buttons := lipgloss.JoinHorizontal(lipgloss.Top, yesButton, "  ", noButton)
	return lipgloss.NewStyle().Align(lipgloss.Right).Render(buttons)
}
