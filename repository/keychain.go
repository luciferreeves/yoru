package repository

import (
	"yoru/database"
	"yoru/models"
)

func GetAllIdentities() ([]models.Identity, error) {
	var identities []models.Identity
	if err := database.DB.Order("id DESC").Find(&identities).Error; err != nil {
		return nil, err
	}
	return identities, nil
}

func GetIdentityByID(id uint) (*models.Identity, error) {
	var identity models.Identity
	if err := database.DB.First(&identity, id).Error; err != nil {
		return nil, err
	}
	return &identity, nil
}

func GetKeyByID(id uint) (*models.Key, error) {
	var key models.Key
	if err := database.DB.First(&key, id).Error; err != nil {
		return nil, err
	}
	return &key, nil
}

func GetAllKeys() ([]models.Key, error) {
	var keys []models.Key
	if err := database.DB.Order("id DESC").Find(&keys).Error; err != nil {
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
