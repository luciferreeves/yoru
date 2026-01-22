package forms

import (
	"net"
	"regexp"
	"strconv"
	"yoru/models"
	"yoru/repository"
	"yoru/screens/styles"
	"yoru/types"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	FieldName int = iota
	FieldHostname
	FieldPort
	FieldMode
	TotalFields
)

const (
	ModeSSH int = iota
	ModeTelnet
	TotalModes
)

type HostForm struct {
	currentHost *models.Host
	activeMode  types.ConnectionMode
	focused     bool
	fieldIndex  int

	nameInput     textinput.Model
	hostnameInput textinput.Model
	portInput     textinput.Model
	modeIndex     int

	fieldErrors        map[int]string
	lastSelectedHostID uint
}

func NewHostForm() *HostForm {
	nameInput := textinput.New()
	nameInput.Placeholder = "Enter host name"
	nameInput.CharLimit = 100
	nameInput.Width = 20
	nameInput.Blur()

	hostnameInput := textinput.New()
	hostnameInput.Placeholder = "192.168.1.1"
	hostnameInput.CharLimit = 15
	hostnameInput.Width = 30
	hostnameInput.Blur()

	portInput := textinput.New()
	portInput.Placeholder = "22"
	portInput.CharLimit = 5
	portInput.Width = 8
	portInput.Blur()

	return &HostForm{
		activeMode:    types.ModeSSH,
		nameInput:     nameInput,
		hostnameInput: hostnameInput,
		portInput:     portInput,
		fieldErrors:   make(map[int]string),
	}
}

func (form *HostForm) LoadHost(host *models.Host) {
	form.currentHost = host
	form.lastSelectedHostID = host.ID
	form.activeMode = host.Mode
	form.fieldIndex = FieldName
	form.modeIndex = ModeSSH
	form.fieldErrors = make(map[int]string)

	if host.Mode == types.ModeTelnet {
		form.modeIndex = ModeTelnet
	}

	form.nameInput.SetValue(host.Name)
	form.hostnameInput.SetValue(host.Hostname)
	form.portInput.SetValue(strconv.Itoa(host.Port))

	form.nameInput.CursorEnd()
	form.hostnameInput.CursorEnd()
	form.portInput.CursorEnd()

	form.setFieldFocus()
}

func (form *HostForm) setFieldFocus() {
	form.nameInput.Blur()
	form.hostnameInput.Blur()
	form.portInput.Blur()

	if !form.focused {
		return
	}

	switch form.fieldIndex {
	case FieldName:
		form.nameInput.Focus()
	case FieldHostname:
		form.hostnameInput.Focus()
	case FieldPort:
		form.portInput.Focus()
	}
}

func (form *HostForm) Save() {
	if form.currentHost != nil {
		form.currentHost.Name = form.nameInput.Value()
		form.currentHost.Hostname = form.hostnameInput.Value()

		if port, err := strconv.Atoi(form.portInput.Value()); err == nil && port > 0 && port <= 65535 {
			form.currentHost.Port = port
		}

		if form.modeIndex == ModeSSH {
			form.currentHost.Mode = types.ModeSSH
		} else {
			form.currentHost.Mode = types.ModeTelnet
		}
		repository.UpdateHost(form.currentHost)
	}
}

func (form *HostForm) validateHostname() {
	value := form.hostnameInput.Value()
	if value == "" {
		form.fieldErrors[FieldHostname] = "Hostname is required"
		return
	}

	if !isValidIP(value) {
		form.fieldErrors[FieldHostname] = "Invalid IP format"
	}
}

func (form *HostForm) validatePort() {
	value := form.portInput.Value()
	if value == "" {
		form.fieldErrors[FieldPort] = "Port is required"
		return
	}

	port, err := strconv.Atoi(value)
	if err != nil {
		form.fieldErrors[FieldPort] = "Port must be a number"
		return
	}

	if port < 1 || port > 65535 {
		form.fieldErrors[FieldPort] = "Port must be 1-65535"
	}
}

func (form *HostForm) validateAll() {
	if form.nameInput.Value() == "" {
		form.fieldErrors[FieldName] = "Name is required"
	}
	form.validateHostname()
	form.validatePort()
}

func (form *HostForm) clearErrors() {
	form.fieldErrors = make(map[int]string)
}

func (form *HostForm) GetError(fieldIndex int) string {
	return form.fieldErrors[fieldIndex]
}

func isValidIP(ip string) bool {
	if net.ParseIP(ip) != nil {
		return true
	}

	ipRegex := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	if !ipRegex.MatchString(ip) {
		return false
	}

	parts := regexp.MustCompile(`\.`).Split(ip, -1)
	for _, part := range parts {
		if num, err := strconv.Atoi(part); err != nil || num > 255 {
			return false
		}
	}

	return true
}

func (form *HostForm) Update(event interface{}) {
	if !form.focused {
		return
	}

	keyMsg, ok := event.(tea.KeyMsg)
	if !ok {
		return
	}

	switch keyMsg.Type {
	case tea.KeyUp:
		if form.fieldIndex > FieldName {
			form.validateCurrentField()
			form.fieldIndex--
			form.setFieldFocus()
		}
		return
	case tea.KeyDown:
		if form.fieldIndex < FieldMode {
			form.validateCurrentField()
			form.fieldIndex++
			form.setFieldFocus()
		}
		return
	case tea.KeySpace:
		if form.fieldIndex == FieldMode {
			form.modeIndex = (form.modeIndex + 1) % TotalModes
			return
		}
	case tea.KeyLeft, tea.KeyRight:
		switch form.fieldIndex {
		case FieldName:
			form.nameInput, _ = form.nameInput.Update(keyMsg)
		case FieldHostname:
			form.hostnameInput, _ = form.hostnameInput.Update(keyMsg)
		case FieldPort:
			form.portInput, _ = form.portInput.Update(keyMsg)
		}
		return
	}

	switch form.fieldIndex {
	case FieldName:
		form.nameInput, _ = form.nameInput.Update(keyMsg)
	case FieldHostname:
		form.hostnameInput, _ = form.hostnameInput.Update(keyMsg)
	case FieldPort:
		if keyMsg.Type == tea.KeyBackspace || keyMsg.Type == tea.KeyDelete ||
			keyMsg.Type == tea.KeyLeft || keyMsg.Type == tea.KeyRight ||
			keyMsg.Type == tea.KeyHome || keyMsg.Type == tea.KeyEnd {
			form.portInput, _ = form.portInput.Update(keyMsg)
		} else if len(keyMsg.Runes) > 0 {
			r := keyMsg.Runes[0]
			if r >= '0' && r <= '9' {
				form.portInput, _ = form.portInput.Update(keyMsg)
			}
		}
	}
}

func (form *HostForm) validateCurrentField() {
	delete(form.fieldErrors, form.fieldIndex)

	switch form.fieldIndex {
	case FieldHostname:
		form.validateHostname()
	case FieldPort:
		form.validatePort()
	}
}

func (form *HostForm) GetLastSelectedHostID() uint {
	return form.lastSelectedHostID
}

func (form *HostForm) SetFocused(focused bool) {
	form.focused = focused
	if focused {
		form.setFieldFocus()
	} else {
		form.nameInput.Blur()
		form.hostnameInput.Blur()
		form.portInput.Blur()
	}
}

func (form *HostForm) Render() string {
	if form.currentHost == nil {
		return styles.FormEmpty.Render("← Select a host or press Ctrl+N to create new")
	}

	return form.renderEditableForm()
}

func (form *HostForm) renderEditableForm() string {
	var lines []string
	lines = append(lines, "")

	nameView := form.nameInput.View()
	nameLine := lipgloss.JoinHorizontal(lipgloss.Left, styles.FormLabel.Render("Name"), nameView)
	if errMsg, ok := form.fieldErrors[FieldName]; ok {
		nameLine = lipgloss.JoinVertical(lipgloss.Left, nameLine, renderError(errMsg))
	}
	lines = append(lines, nameLine)

	hostnameView := form.hostnameInput.View()
	hostnameLine := lipgloss.JoinHorizontal(lipgloss.Left, styles.FormLabel.Render("Hostname"), hostnameView)
	if errMsg, ok := form.fieldErrors[FieldHostname]; ok {
		hostnameLine = lipgloss.JoinVertical(lipgloss.Left, hostnameLine, renderError(errMsg))
	}
	lines = append(lines, hostnameLine)

	portView := form.portInput.View()
	portLine := lipgloss.JoinHorizontal(lipgloss.Left, styles.FormLabel.Render("Port"), portView)
	if errMsg, ok := form.fieldErrors[FieldPort]; ok {
		portLine = lipgloss.JoinVertical(lipgloss.Left, portLine, renderError(errMsg))
	}
	lines = append(lines, portLine)

	modeView := form.renderModeChooser()
	lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Left, styles.FormLabel.Render("Mode"), modeView))

	lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Left, styles.FormLabel.Render("Identity"), styles.FormPlaceholder.Render("(Coming soon)")))
	lines = append(lines, "")

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func renderError(errMsg string) string {
	return styles.FormError.Render("✗ " + errMsg)
}

func (form *HostForm) renderModeChooser() string {
	sshBox := "[ ]"
	telnetBox := "[ ]"

	if form.modeIndex == ModeSSH {
		sshBox = "[x]"
	} else {
		telnetBox = "[x]"
	}

	style := styles.FormText
	if form.focused && form.fieldIndex == FieldMode {
		style = styles.FormTextFocused
	}

	sshPart := style.Render(sshBox + " SSH")
	telnetPart := style.Render(telnetBox + " Telnet")

	return sshPart + "     " + telnetPart
}
