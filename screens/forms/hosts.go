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
	FieldIdentity
	TotalFields
)

const (
	ModeSSH int = iota
	ModeTelnet
	TotalModes
)

type HostForm struct {
	currentHost      *models.Host
	activeMode       types.ConnectionMode
	focused          bool
	fieldIndex       int
	selectedCredType types.CredentialType
	selectedCredID   uint

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

	if host.CredentialID > 0 {
		form.selectedCredType = host.CredentialType
		form.selectedCredID = host.CredentialID
	} else {
		form.selectedCredType = ""
		form.selectedCredID = 0
	}

	form.nameInput.SetValue(host.Name)
	form.hostnameInput.SetValue(host.Hostname)
	form.portInput.SetValue(strconv.Itoa(host.Port))

	form.nameInput.CursorEnd()
	form.hostnameInput.CursorEnd()
	form.portInput.CursorEnd()

	form.setFieldFocus()
}

func (form *HostForm) Clear() {
	form.currentHost = nil
	form.lastSelectedHostID = 0
	form.selectedCredType = ""
	form.selectedCredID = 0
	form.fieldErrors = make(map[int]string)
	form.nameInput.SetValue("")
	form.hostnameInput.SetValue("")
	form.portInput.SetValue("")
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

		if form.selectedCredID > 0 {
			form.currentHost.CredentialType = form.selectedCredType
			form.currentHost.CredentialID = form.selectedCredID
		} else {
			form.currentHost.CredentialID = 0
			form.currentHost.CredentialType = ""
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
		if form.fieldIndex < FieldIdentity {
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
		emptyMsg := styles.FormEmpty.Render("← Nothing to see here! Press Ctrl+N to add a new host first.")
		return lipgloss.Place(
			lipgloss.Width(emptyMsg)+4,
			lipgloss.Height(emptyMsg)+4,
			lipgloss.Center,
			lipgloss.Center,
			emptyMsg,
		)
	}

	return form.renderEditableForm()
}

func (form *HostForm) renderEditableForm() string {
	var fields []string

	fields = append(fields, styles.FormSectionTitle.Render("Connection Details"))

	var nameLabel string
	if form.focused && form.fieldIndex == FieldName {
		nameLabel = styles.FormLabelFocused.Render("Name")
	} else {
		nameLabel = styles.FormLabel.Render("Name")
	}
	var nameView string
	if form.focused && form.fieldIndex == FieldName {
		nameView = styles.FormInputFocused.Render(form.nameInput.View())
	} else {
		nameView = styles.FormInput.Render(form.nameInput.View())
	}
	nameLine := lipgloss.JoinHorizontal(lipgloss.Left, nameLabel, nameView)
	fields = append(fields, styles.FormFieldContainer.Render(nameLine))
	if errMsg, ok := form.fieldErrors[FieldName]; ok {
		fields = append(fields, renderError(errMsg))
	}

	var hostnameLabel string
	if form.focused && form.fieldIndex == FieldHostname {
		hostnameLabel = styles.FormLabelFocused.Render("Hostname")
	} else {
		hostnameLabel = styles.FormLabel.Render("Hostname")
	}
	var hostnameView string
	if form.focused && form.fieldIndex == FieldHostname {
		hostnameView = styles.FormInputFocused.Render(form.hostnameInput.View())
	} else {
		hostnameView = styles.FormInput.Render(form.hostnameInput.View())
	}
	hostnameLine := lipgloss.JoinHorizontal(lipgloss.Left, hostnameLabel, hostnameView)
	fields = append(fields, styles.FormFieldContainer.Render(hostnameLine))
	if errMsg, ok := form.fieldErrors[FieldHostname]; ok {
		fields = append(fields, renderError(errMsg))
	}

	var portLabel string
	if form.focused && form.fieldIndex == FieldPort {
		portLabel = styles.FormLabelFocused.Render("Port")
	} else {
		portLabel = styles.FormLabel.Render("Port")
	}
	var portView string
	if form.focused && form.fieldIndex == FieldPort {
		portView = styles.FormInputFocused.Render(form.portInput.View())
	} else {
		portView = styles.FormInput.Render(form.portInput.View())
	}
	portLine := lipgloss.JoinHorizontal(lipgloss.Left, portLabel, portView)
	fields = append(fields, styles.FormFieldContainer.Render(portLine))
	if errMsg, ok := form.fieldErrors[FieldPort]; ok {
		fields = append(fields, renderError(errMsg))
	}

	var modeLabel string
	if form.focused && form.fieldIndex == FieldMode {
		modeLabel = styles.FormLabelFocused.Render("Mode")
	} else {
		modeLabel = styles.FormLabel.Render("Mode")
	}
	modeView := form.renderModeChooser()
	modeLine := lipgloss.JoinHorizontal(lipgloss.Left, modeLabel, modeView)
	fields = append(fields, styles.FormFieldContainer.Render(modeLine))

	fields = append(fields, styles.FormSectionTitle.Render("Authentication"))

	var identityLabel string
	if form.focused && form.fieldIndex == FieldIdentity {
		identityLabel = styles.FormLabelFocused.Render("Identity")
	} else {
		identityLabel = styles.FormLabel.Render("Identity")
	}

	var identityText string
	var identityStyle lipgloss.Style
	if form.selectedCredID > 0 {
		// Fetch the credential name based on type
		switch form.selectedCredType {
		case types.CredentialIdentity:
			if identity, err := repository.GetIdentityByID(form.selectedCredID); err == nil {
				identityText = identity.Name
			} else {
				identityText = "Unknown identity"
			}
		case types.CredentialKey:
			if key, err := repository.GetKeyByID(form.selectedCredID); err == nil {
				identityText = key.Name + " (Key)"
			} else {
				identityText = "Unknown key"
			}
		default:
			identityText = "Unknown credential"
		}
		if form.focused && form.fieldIndex == FieldIdentity {
			identityStyle = styles.FormInputFocused
		} else {
			identityStyle = styles.FormInput
		}
	} else {
		identityText = "No identity selected"
		if form.focused && form.fieldIndex == FieldIdentity {
			identityStyle = styles.FormInputFocused
		} else {
			identityStyle = styles.FormPlaceholder
		}
	}

	identityView := identityStyle.Render(identityText)
	identityLine := lipgloss.JoinHorizontal(lipgloss.Left, identityLabel, identityView)
	fields = append(fields, styles.FormFieldContainer.Render(identityLine))

	formContent := lipgloss.JoinVertical(lipgloss.Left, fields...)
	return styles.FormContainer.Render(formContent)
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

	checkboxStyle := styles.FormCheckbox
	labelStyle := styles.FormCheckboxLabel
	if form.focused && form.fieldIndex == FieldMode {
		checkboxStyle = styles.FormCheckboxFocused
		labelStyle = styles.FormCheckboxLabelFocused
	}

	sshPart := lipgloss.JoinHorizontal(lipgloss.Left,
		checkboxStyle.Render(sshBox),
		labelStyle.Render("SSH"))
	telnetPart := lipgloss.JoinHorizontal(lipgloss.Left,
		checkboxStyle.Render(telnetBox),
		labelStyle.Render("Telnet"))

	return lipgloss.JoinHorizontal(lipgloss.Left, sshPart, "   ", telnetPart)
}

func (form *HostForm) GetFieldIndex() int {
	return form.fieldIndex
}

func (form *HostForm) SetSelectedCredential(credType types.CredentialType, credID uint) {
	form.selectedCredType = credType
	form.selectedCredID = credID
}

func (form *HostForm) GetSelectedCredential() (types.CredentialType, uint) {
	return form.selectedCredType, form.selectedCredID
}
