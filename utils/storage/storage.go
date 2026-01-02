package storage

import (
	"yoru/types"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDatabaseInstance() (*types.Database, error) {
	databasePath, err := getDatabasePath()
	if err != nil {
		return nil, err
	}

	dialector := sqlite.Open(databasePath)
	database, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	return &types.Database{DB: database}, err
}
