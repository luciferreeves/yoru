package forms

import (
	"yoru/models"
	"yoru/repository"
	"yoru/screens/components"
	"yoru/screens/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	KeychainFieldType int = iota
	KeychainFieldName
	KeychainFieldUsername
	KeychainFieldPassword
	KeychainFieldPrivateKey
	KeychainFieldPublicKey
	KeychainFieldCertificate
	TotalKeychainFields
)

type KeychainForm struct {
	currentKey      *models.Key
	currentIdentity *models.Identity
	itemType        string // "Key" or "Identity"
	focused         bool
	fieldIndex      int
	typeIndex       int // 0 = Identity, 1 = Key

	nameInput     textinput.Model
	usernameInput textinput.Model
	passwordInput textinput.Model

	privateKeyArea  components.TextArea
	publicKeyArea   components.TextArea
	certificateArea components.TextArea

	fieldErrors      map[int]string
	lastSelectedID   uint
	lastSelectedType string
}

func NewKeychainForm() *KeychainForm {
	nameInput := textinput.New()
	nameInput.Placeholder = "Enter name"
	nameInput.CharLimit = 100
	nameInput.Width = 40
	nameInput.Blur()

	usernameInput := textinput.New()
	usernameInput.Placeholder = "Enter username"
	usernameInput.CharLimit = 100
	usernameInput.Width = 40
	usernameInput.Blur()

	passwordInput := textinput.New()
	passwordInput.Placeholder = "Enter password"
	passwordInput.CharLimit = 100
	passwordInput.Width = 40
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.Blur()

	privateKeyArea := components.NewTextArea("-----BEGIN OPENSSH PRIVATE KEY-----", 60, 5)
	publicKeyArea := components.NewTextArea("ssh-rsa AAAA...", 60, 3)
	certificateArea := components.NewTextArea("-----BEGIN CERTIFICATE-----", 60, 3)

	return &KeychainForm{
		nameInput:        nameInput,
		usernameInput:    usernameInput,
		passwordInput:    passwordInput,
		privateKeyArea:   privateKeyArea,
		publicKeyArea:    publicKeyArea,
		certificateArea:  certificateArea,
		fieldErrors:      make(map[int]string),
		lastSelectedID:   0,
		lastSelectedType: "",
	}
}

func (form *KeychainForm) LoadKey(key *models.Key) {
	form.currentKey = key
	form.currentIdentity = nil
	form.itemType = "Key"
	form.typeIndex = 1
	form.lastSelectedID = key.ID
	form.lastSelectedType = "Key"
	form.fieldIndex = KeychainFieldName
	form.fieldErrors = make(map[int]string)

	form.nameInput.SetValue(key.Name)
	form.usernameInput.SetValue(key.Username)
	form.privateKeyArea.SetValue(key.PrivateKey)
	form.publicKeyArea.SetValue(key.PublicKey)
	form.certificateArea.SetValue(key.Certificate)

	form.nameInput.CursorEnd()
	form.setFieldFocus()
}

func (form *KeychainForm) LoadIdentity(identity *models.Identity) {
	form.currentIdentity = identity
	form.currentKey = nil
	form.itemType = "Identity"
	form.typeIndex = 0
	form.lastSelectedID = identity.ID
	form.lastSelectedType = "Identity"
	form.fieldIndex = KeychainFieldName
	form.fieldErrors = make(map[int]string)

	form.nameInput.SetValue(identity.Name)
	form.usernameInput.SetValue(identity.Username)
	form.passwordInput.SetValue(identity.Password)

	form.nameInput.CursorEnd()
	form.setFieldFocus()
}

func (form *KeychainForm) Clear() {
	form.currentKey = nil
	form.currentIdentity = nil
	form.lastSelectedID = 0
	form.lastSelectedType = ""
	form.fieldErrors = make(map[int]string)
	form.nameInput.SetValue("")
	form.usernameInput.SetValue("")
	form.passwordInput.SetValue("")
	form.privateKeyArea.Reset()
	form.publicKeyArea.Reset()
	form.certificateArea.Reset()
}

func (form *KeychainForm) setFieldFocus() {
	form.nameInput.Blur()
	form.usernameInput.Blur()
	form.passwordInput.Blur()
	form.privateKeyArea.Blur()
	form.publicKeyArea.Blur()
	form.certificateArea.Blur()

	if !form.focused {
		return
	}

	if form.itemType == "Key" {
		switch form.fieldIndex {
		case KeychainFieldName:
			form.nameInput.Focus()
		case KeychainFieldUsername:
			form.usernameInput.Focus()
		case KeychainFieldPrivateKey:
			form.privateKeyArea.Focus()
		case KeychainFieldPublicKey:
			form.publicKeyArea.Focus()
		case KeychainFieldCertificate:
			form.certificateArea.Focus()
		}
	} else {
		switch form.fieldIndex {
		case KeychainFieldName:
			form.nameInput.Focus()
		case KeychainFieldUsername:
			form.usernameInput.Focus()
		case KeychainFieldPassword:
			form.passwordInput.Focus()
		}
	}
}

func (form *KeychainForm) validateKey() {
	if form.nameInput.Value() == "" {
		form.fieldErrors[KeychainFieldName] = "Name is required"
	}
	if form.usernameInput.Value() == "" {
		form.fieldErrors[KeychainFieldUsername] = "Username is required"
	}
	if form.privateKeyArea.Value() == "" {
		form.fieldErrors[KeychainFieldPrivateKey] = "Private key is required"
	}
}

func (form *KeychainForm) validateIdentity() {
	if form.nameInput.Value() == "" {
		form.fieldErrors[KeychainFieldName] = "Name is required"
	}
	if form.usernameInput.Value() == "" {
		form.fieldErrors[KeychainFieldUsername] = "Username is required"
	}
	if form.passwordInput.Value() == "" {
		form.fieldErrors[KeychainFieldPassword] = "Password is required"
	}
}

func (form *KeychainForm) Save() {
	if form.itemType == "Key" && form.currentKey != nil {
		form.validateKey()
		if len(form.fieldErrors) > 0 {
			return
		}

		form.currentKey.Name = form.nameInput.Value()
		form.currentKey.Username = form.usernameInput.Value()
		form.currentKey.PrivateKey = form.privateKeyArea.Value()
		form.currentKey.PublicKey = form.publicKeyArea.Value()
		form.currentKey.Certificate = form.certificateArea.Value()
		repository.UpdateKey(form.currentKey)
	} else if form.itemType == "Identity" && form.currentIdentity != nil {
		form.validateIdentity()
		if len(form.fieldErrors) > 0 {
			return
		}

		form.currentIdentity.Name = form.nameInput.Value()
		form.currentIdentity.Username = form.usernameInput.Value()
		form.currentIdentity.Password = form.passwordInput.Value()
		repository.UpdateIdentity(form.currentIdentity)
	}
}

func (form *KeychainForm) GetLastSelectedID() uint {
	return form.lastSelectedID
}

func (form *KeychainForm) GetLastSelectedType() string {
	return form.lastSelectedType
}

func (form *KeychainForm) IsTextAreaEditing() bool {
	return form.privateKeyArea.IsEditing() || form.publicKeyArea.IsEditing() || form.certificateArea.IsEditing()
}

func (form *KeychainForm) SetFocused(focused bool) {
	form.focused = focused
	if focused {
		form.fieldIndex = KeychainFieldType
		form.setFieldFocus()
	} else {
		form.nameInput.Blur()
		form.usernameInput.Blur()
		form.passwordInput.Blur()
		form.privateKeyArea.Blur()
		form.publicKeyArea.Blur()
		form.certificateArea.Blur()
	}
}

func (form *KeychainForm) GetError(fieldIndex int) string {
	return form.fieldErrors[fieldIndex]
}

func (form *KeychainForm) Update(event interface{}) {
	if !form.focused {
		return
	}

	keyMsg, ok := event.(tea.KeyMsg)
	if !ok {
		return
	}

	maxFieldIndex := 0
	if form.itemType == "Key" {
		maxFieldIndex = KeychainFieldCertificate
	} else {
		maxFieldIndex = KeychainFieldPassword
	}

	isEditing := form.IsTextAreaEditing()

	switch keyMsg.Type {
	case tea.KeyEnter:
		if !isEditing && form.itemType == "Key" {
			switch form.fieldIndex {
			case KeychainFieldPrivateKey:
				form.privateKeyArea.StartEditing()
				return
			case KeychainFieldPublicKey:
				form.publicKeyArea.StartEditing()
				return
			case KeychainFieldCertificate:
				form.certificateArea.StartEditing()
				return
			}
		}
	case tea.KeyEscape:
		if isEditing {
			switch form.fieldIndex {
			case KeychainFieldPrivateKey:
				form.privateKeyArea.StopEditing()
			case KeychainFieldPublicKey:
				form.publicKeyArea.StopEditing()
			case KeychainFieldCertificate:
				form.certificateArea.StopEditing()
			}
			return
		}
	case tea.KeyUp:
		if isEditing {
		} else if form.fieldIndex > KeychainFieldType {
			delete(form.fieldErrors, form.fieldIndex)
			form.fieldIndex--
			if form.itemType == "Key" && form.fieldIndex == KeychainFieldPassword {
				form.fieldIndex = KeychainFieldUsername
			}
			if form.fieldIndex < KeychainFieldType {
				form.fieldIndex = KeychainFieldType
			}
			form.setFieldFocus()
			return
		}
	case tea.KeyDown:
		if isEditing {
		} else if form.fieldIndex < maxFieldIndex {
			delete(form.fieldErrors, form.fieldIndex)
			form.fieldIndex++
			if form.itemType == "Key" && form.fieldIndex == KeychainFieldPassword {
				form.fieldIndex = KeychainFieldPrivateKey
			}
			if form.fieldIndex > maxFieldIndex {
				form.fieldIndex = maxFieldIndex
			}
			form.setFieldFocus()
			return
		}
	case tea.KeySpace:
		if form.fieldIndex == KeychainFieldType {
			form.typeIndex = (form.typeIndex + 1) % 2
			if form.typeIndex == 0 {
				form.switchToIdentity()
			} else {
				form.switchToKey()
			}
			return
		}
	}

	if form.itemType == "Key" {
		switch form.fieldIndex {
		case KeychainFieldName:
			form.nameInput, _ = form.nameInput.Update(keyMsg)
		case KeychainFieldUsername:
			form.usernameInput, _ = form.usernameInput.Update(keyMsg)
		case KeychainFieldPrivateKey:
			form.privateKeyArea.Update(keyMsg)
		case KeychainFieldPublicKey:
			form.publicKeyArea.Update(keyMsg)
		case KeychainFieldCertificate:
			form.certificateArea.Update(keyMsg)
		}
	} else {
		switch form.fieldIndex {
		case KeychainFieldName:
			form.nameInput, _ = form.nameInput.Update(keyMsg)
		case KeychainFieldUsername:
			form.usernameInput, _ = form.usernameInput.Update(keyMsg)
		case KeychainFieldPassword:
			form.passwordInput, _ = form.passwordInput.Update(keyMsg)
		}
	}
}

func (form *KeychainForm) Render() string {
	if form.currentKey == nil && form.currentIdentity == nil {
		emptyMsg := styles.FormEmpty.Render("← Select an item or press Ctrl+N to create new")
		return lipgloss.Place(
			lipgloss.Width(emptyMsg)+4,
			lipgloss.Height(emptyMsg)+4,
			lipgloss.Center,
			lipgloss.Center,
			emptyMsg,
		)
	}

	if form.itemType == "Key" {
		return form.renderKeyForm()
	}
	return form.renderIdentityForm()
}

func (form *KeychainForm) renderKeyForm() string {
	var fields []string

	var typeLabel string
	if form.focused && form.fieldIndex == KeychainFieldType {
		typeLabel = styles.FormLabelFocused.Render("Type")
	} else {
		typeLabel = styles.FormLabel.Render("Type")
	}
	typeView := form.renderTypeChooser()
	typeLine := lipgloss.JoinHorizontal(lipgloss.Left, typeLabel, typeView)
	fields = append(fields, styles.FormFieldContainer.Render(typeLine))

	fields = append(fields, styles.FormSectionTitle.Render("Connection Details"))

	var nameLabel string
	if form.focused && form.fieldIndex == KeychainFieldName {
		nameLabel = styles.FormLabelFocused.Render("Name")
	} else {
		nameLabel = styles.FormLabel.Render("Name")
	}
	var nameView string
	if form.focused && form.fieldIndex == KeychainFieldName {
		nameView = styles.FormInputFocused.Render(form.nameInput.View())
	} else {
		nameView = styles.FormInput.Render(form.nameInput.View())
	}
	nameLine := lipgloss.JoinHorizontal(lipgloss.Left, nameLabel, nameView)
	fields = append(fields, styles.FormFieldContainer.Render(nameLine))
	if errMsg, ok := form.fieldErrors[KeychainFieldName]; ok {
		fields = append(fields, renderKeychainError(errMsg))
	}

	var usernameLabel string
	if form.focused && form.fieldIndex == KeychainFieldUsername {
		usernameLabel = styles.FormLabelFocused.Render("Username")
	} else {
		usernameLabel = styles.FormLabel.Render("Username")
	}
	var usernameView string
	if form.focused && form.fieldIndex == KeychainFieldUsername {
		usernameView = styles.FormInputFocused.Render(form.usernameInput.View())
	} else {
		usernameView = styles.FormInput.Render(form.usernameInput.View())
	}
	usernameLine := lipgloss.JoinHorizontal(lipgloss.Left, usernameLabel, usernameView)
	fields = append(fields, styles.FormFieldContainer.Render(usernameLine))
	if errMsg, ok := form.fieldErrors[KeychainFieldUsername]; ok {
		fields = append(fields, renderKeychainError(errMsg))
	}

	var privateKeyLabel string
	if form.focused && form.fieldIndex == KeychainFieldPrivateKey {
		privateKeyLabel = styles.FormLabelFocused.Render("Private Key")
	} else {
		privateKeyLabel = styles.FormLabel.Render("Private Key")
	}
	var privateKeyView string
	if form.focused && form.fieldIndex == KeychainFieldPrivateKey {
		if form.privateKeyArea.IsEditing() {
			privateKeyView = styles.FormTextAreaEditing.Render(form.privateKeyArea.View())
		} else {
			privateKeyView = styles.FormTextAreaFocused.Render(form.privateKeyArea.View())
		}
	} else {
		privateKeyView = styles.FormTextAreaInvisibleBorder.Render(form.privateKeyArea.View())
	}
	privateKeyLine := lipgloss.JoinHorizontal(lipgloss.Left, privateKeyLabel, privateKeyView)
	fields = append(fields, styles.FormFieldContainer.Render(privateKeyLine))
	if errMsg, ok := form.fieldErrors[KeychainFieldPrivateKey]; ok {
		fields = append(fields, renderKeychainError(errMsg))
	}

	var publicKeyLabel string
	if form.focused && form.fieldIndex == KeychainFieldPublicKey {
		publicKeyLabel = styles.FormLabelFocused.Render("Public Key")
	} else {
		publicKeyLabel = styles.FormLabel.Render("Public Key")
	}
	var publicKeyView string
	if form.focused && form.fieldIndex == KeychainFieldPublicKey {
		if form.publicKeyArea.IsEditing() {
			publicKeyView = styles.FormTextAreaEditing.Render(form.publicKeyArea.View())
		} else {
			publicKeyView = styles.FormTextAreaFocused.Render(form.publicKeyArea.View())
		}
	} else {
		publicKeyView = styles.FormTextAreaInvisibleBorder.Render(form.publicKeyArea.View())
	}
	publicKeyLine := lipgloss.JoinHorizontal(lipgloss.Left, publicKeyLabel, publicKeyView)
	fields = append(fields, styles.FormFieldContainer.Render(publicKeyLine))

	var certificateLabel string
	if form.focused && form.fieldIndex == KeychainFieldCertificate {
		certificateLabel = styles.FormLabelFocused.Render("Certificate")
	} else {
		certificateLabel = styles.FormLabel.Render("Certificate")
	}
	var certificateView string
	if form.focused && form.fieldIndex == KeychainFieldCertificate {
		if form.certificateArea.IsEditing() {
			certificateView = styles.FormTextAreaEditing.Render(form.certificateArea.View())
		} else {
			certificateView = styles.FormTextAreaFocused.Render(form.certificateArea.View())
		}
	} else {
		certificateView = styles.FormTextAreaInvisibleBorder.Render(form.certificateArea.View())
	}
	certificateLine := lipgloss.JoinHorizontal(lipgloss.Left, certificateLabel, certificateView)
	fields = append(fields, styles.FormFieldContainer.Render(certificateLine))

	formContent := lipgloss.JoinVertical(lipgloss.Left, fields...)
	return styles.FormContainer.Render(formContent)
}

func (form *KeychainForm) renderIdentityForm() string {
	var fields []string

	var typeLabel string
	if form.focused && form.fieldIndex == KeychainFieldType {
		typeLabel = styles.FormLabelFocused.Render("Type")
	} else {
		typeLabel = styles.FormLabel.Render("Type")
	}
	typeView := form.renderTypeChooser()
	typeLine := lipgloss.JoinHorizontal(lipgloss.Left, typeLabel, typeView)
	fields = append(fields, styles.FormFieldContainer.Render(typeLine))

	fields = append(fields, styles.FormSectionTitle.Render("Connection Details"))

	var nameLabel string
	if form.focused && form.fieldIndex == KeychainFieldName {
		nameLabel = styles.FormLabelFocused.Render("Name")
	} else {
		nameLabel = styles.FormLabel.Render("Name")
	}
	var nameView string
	if form.focused && form.fieldIndex == KeychainFieldName {
		nameView = styles.FormInputFocused.Render(form.nameInput.View())
	} else {
		nameView = styles.FormInput.Render(form.nameInput.View())
	}
	nameLine := lipgloss.JoinHorizontal(lipgloss.Left, nameLabel, nameView)
	fields = append(fields, styles.FormFieldContainer.Render(nameLine))
	if errMsg, ok := form.fieldErrors[KeychainFieldName]; ok {
		fields = append(fields, renderKeychainError(errMsg))
	}

	var usernameLabel string
	if form.focused && form.fieldIndex == KeychainFieldUsername {
		usernameLabel = styles.FormLabelFocused.Render("Username")
	} else {
		usernameLabel = styles.FormLabel.Render("Username")
	}
	var usernameView string
	if form.focused && form.fieldIndex == KeychainFieldUsername {
		usernameView = styles.FormInputFocused.Render(form.usernameInput.View())
	} else {
		usernameView = styles.FormInput.Render(form.usernameInput.View())
	}
	usernameLine := lipgloss.JoinHorizontal(lipgloss.Left, usernameLabel, usernameView)
	fields = append(fields, styles.FormFieldContainer.Render(usernameLine))
	if errMsg, ok := form.fieldErrors[KeychainFieldUsername]; ok {
		fields = append(fields, renderKeychainError(errMsg))
	}

	var passwordLabel string
	if form.focused && form.fieldIndex == KeychainFieldPassword {
		passwordLabel = styles.FormLabelFocused.Render("Password")
	} else {
		passwordLabel = styles.FormLabel.Render("Password")
	}
	var passwordView string
	if form.focused && form.fieldIndex == KeychainFieldPassword {
		passwordView = styles.FormInputFocused.Render(form.passwordInput.View())
	} else {
		passwordView = styles.FormInput.Render(form.passwordInput.View())
	}
	passwordLine := lipgloss.JoinHorizontal(lipgloss.Left, passwordLabel, passwordView)
	fields = append(fields, styles.FormFieldContainer.Render(passwordLine))
	if errMsg, ok := form.fieldErrors[KeychainFieldPassword]; ok {
		fields = append(fields, renderKeychainError(errMsg))
	}

	formContent := lipgloss.JoinVertical(lipgloss.Left, fields...)
	return styles.FormContainer.Render(formContent)
}

func renderKeychainError(errMsg string) string {
	return styles.FormError.Render("✗ " + errMsg)
}

func (form *KeychainForm) switchToIdentity() {
	name := form.nameInput.Value()

	if form.currentKey != nil {
		repository.DeleteKey(form.currentKey.ID)
		newIdentity := &models.Identity{
			Name:     name,
			Username: "",
			Password: "",
		}
		repository.CreateIdentity(newIdentity)
		form.currentIdentity = newIdentity
		form.currentKey = nil
	}

	form.itemType = "Identity"
	form.usernameInput.SetValue("")
	form.passwordInput.SetValue("")
	form.fieldIndex = KeychainFieldType
	form.setFieldFocus()
}

func (form *KeychainForm) switchToKey() {
	name := form.nameInput.Value()

	if form.currentIdentity != nil {
		repository.DeleteIdentity(form.currentIdentity.ID)
		newKey := &models.Key{
			Name:        name,
			PrivateKey:  "",
			PublicKey:   "",
			Certificate: "",
		}
		repository.CreateKey(newKey)
		form.currentKey = newKey
		form.currentIdentity = nil
	}

	form.itemType = "Key"
	form.privateKeyArea.SetValue("")
	form.publicKeyArea.SetValue("")
	form.certificateArea.SetValue("")
	form.fieldIndex = KeychainFieldType
	form.setFieldFocus()
}

func (form *KeychainForm) renderTypeChooser() string {
	identityBox := "[ ]"
	keyBox := "[ ]"

	if form.typeIndex == 0 {
		identityBox = "[x]"
	} else {
		keyBox = "[x]"
	}

	checkboxStyle := styles.FormCheckbox
	labelStyle := styles.FormCheckboxLabel
	if form.focused && form.fieldIndex == KeychainFieldType {
		checkboxStyle = styles.FormCheckboxFocused
		labelStyle = styles.FormCheckboxLabelFocused
	}

	identityPart := lipgloss.JoinHorizontal(lipgloss.Left,
		checkboxStyle.Render(identityBox),
		labelStyle.Render("Identity"))
	keyPart := lipgloss.JoinHorizontal(lipgloss.Left,
		checkboxStyle.Render(keyBox),
		labelStyle.Render("Key"))

	return lipgloss.JoinHorizontal(lipgloss.Left, identityPart, "   ", keyPart)
}
