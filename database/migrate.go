package database

import (
	"yoru/models"
	"yoru/types"
)

func migrate(db *types.Database) error {
	return db.AutoMigrate(
		&models.Host{},
		&models.KnownHost{},
		&models.Identity{},
		&models.Key{},
		&models.ConnectionLog{},
	)
}
