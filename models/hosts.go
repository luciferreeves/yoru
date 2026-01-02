package models

import (
	"time"
	"yoru/types"
)

type Host struct {
	types.Model
	Name            string               `gorm:"not null"`
	Hostname        string               `gorm:"not null"`
	Mode            types.ConnectionMode `gorm:"type:text;not null"`
	Port            int                  `gorm:"not null"`
	CredentialID    uint                 `gorm:"not null"`
	CredentialType  types.CredentialType `gorm:"type:text;not null"`
	LastConnectedAt *time.Time
}

type KnownHost struct {
	types.Model
	Hostname    string `gorm:"not null"`
	Port        int    `gorm:"not null"`
	KeyType     string `gorm:"not null"`
	Fingerprint string `gorm:"not null;uniqueIndex"`
}
