package repository

import (
	"yoru/database"
	"yoru/models"
)

func GetAllIdentities() ([]models.Identity, error) {
	var identities []models.Identity
	if err := database.DB.Order("created_at DESC").Find(&identities).Error; err != nil {
		return nil, err
	}
	return identities, nil
}

func GetAllKeys() ([]models.Key, error) {
	var keys []models.Key
	if err := database.DB.Order("created_at DESC").Find(&keys).Error; err != nil {
		return nil, err
	}
	return keys, nil
}

func CreateIdentity(identity *models.Identity) error {
	return database.DB.Create(identity).Error
}

func CreateKey(key *models.Key) error {
	return database.DB.Create(key).Error
}

func UpdateIdentity(identity *models.Identity) error {
	return database.DB.Save(identity).Error
}

func UpdateKey(key *models.Key) error {
	return database.DB.Save(key).Error
}

func DeleteIdentity(id uint) error {
	return database.DB.Delete(&models.Identity{}, id).Error
}

func DeleteKey(id uint) error {
	return database.DB.Delete(&models.Key{}, id).Error
}
