package repository

import (
	"yoru/database"
	"yoru/models"
)

func CreateConnectionLog(log *models.ConnectionLog) error {
	return database.DB.Create(log).Error
}

func GetLastNConnectionLogs(limit int) ([]models.ConnectionLog, error) {
	var logs []models.ConnectionLog
	err := database.DB.Order("started_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}

func GetAllConnectionLogs() ([]models.ConnectionLog, error) {
	var logs []models.ConnectionLog
	err := database.DB.Order("started_at DESC").Find(&logs).Error
	return logs, err
}

func GetConnectionLogByID(id uint) (*models.ConnectionLog, error) {
	var log models.ConnectionLog
	err := database.DB.First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func UpdateConnectionLog(log *models.ConnectionLog) error {
	return database.DB.Save(log).Error
}
