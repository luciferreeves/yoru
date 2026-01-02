package database

import (
	"yoru/types"
	"yoru/utils/errors"
	"yoru/utils/storage"
)

var DB *types.Database

func init() {
	var err error
	DB, err = storage.GetDatabaseInstance()
	if err != nil {
		errors.ExitOnDatabaseConnectionFailed(err)
	}

	if err := migrate(DB); err != nil {
		errors.ExitOnMigrationFailed(err)
	}
}
