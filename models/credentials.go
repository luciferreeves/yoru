package models

import "yoru/types"

type Identity struct {
	types.Model
	Name     string `gorm:"not null"`
	Username string `gorm:"not null"`
	Password string `gorm:"not null"`
}

type Key struct {
	types.Model
	Name        string `gorm:"not null"`
	Username    string `gorm:"not null;default:''"`
	PrivateKey  string `gorm:"type:text;not null"`
	PublicKey   string `gorm:"type:text"`
	Certificate string `gorm:"type:text"`
}
