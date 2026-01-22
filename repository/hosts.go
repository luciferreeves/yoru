package repository

import (
	"yoru/database"
	"yoru/models"
	"yoru/types"
)

func CreateHost(host *models.Host) error {
	return database.DB.Create(host).Error
}

func GetAllHosts() ([]models.Host, error) {
	var hosts []models.Host
	err := database.DB.Order("id DESC").Find(&hosts).Error
	return hosts, err
}

func GetHostByID(id uint) (*models.Host, error) {
	var host models.Host
	err := database.DB.First(&host, id).Error
	if err != nil {
		return nil, err
	}
	return &host, nil
}

func GetHostsByMode(mode types.ConnectionMode) ([]models.Host, error) {
	var hosts []models.Host
	err := database.DB.Where("mode = ?", mode).Find(&hosts).Error
	return hosts, err
}

func UpdateHost(host *models.Host) error {
	return database.DB.Save(host).Error
}

func DeleteHost(id uint) error {
	return database.DB.Delete(&models.Host{}, id).Error
}

func UpdateLastConnected(id uint) error {
	return database.DB.Model(&models.Host{}).Where("id = ?", id).Update("last_connected_at", "now()").Error
}

func CreateKnownHost(knownHost *models.KnownHost) error {
	return database.DB.Create(knownHost).Error
}

func GetAllKnownHosts() ([]models.KnownHost, error) {
	var knownHosts []models.KnownHost
	err := database.DB.Find(&knownHosts).Error
	return knownHosts, err
}

func GetKnownHostByFingerprint(fingerprint string) (*models.KnownHost, error) {
	var knownHost models.KnownHost
	err := database.DB.Where("fingerprint = ?", fingerprint).First(&knownHost).Error
	if err != nil {
		return nil, err
	}
	return &knownHost, nil
}

func DeleteKnownHost(id uint) error {
	return database.DB.Delete(&models.KnownHost{}, id).Error
}
