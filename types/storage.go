package types

import "gorm.io/gorm"

type Database struct {
	*gorm.DB
}

type Model struct {
	gorm.Model
}

type ConnectionMode string

const (
	ModeSSH    ConnectionMode = "ssh"
	ModeTelnet ConnectionMode = "telnet"
)

type CredentialType string

const (
	CredentialIdentity CredentialType = "identity"
	CredentialKey      CredentialType = "key"
)
