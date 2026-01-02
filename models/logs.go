package models

import (
	"time"
	"yoru/types"
)

type ConnectionLog struct {
	types.Model
	StartedAt      time.Time `gorm:"not null"`
	EndedAt        *time.Time
	LocalHostname  string               `gorm:"not null"`
	LocalIP        string               `gorm:"not null"`
	RemoteHostname string               `gorm:"not null"`
	Mode           types.ConnectionMode `gorm:"type:text;not null"`
	CredentialID   uint                 `gorm:"not null"`
	CredentialType types.CredentialType `gorm:"type:text;not null"`
}
