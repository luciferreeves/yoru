package errors

import (
	"fmt"
	"os"
	"yoru/shared"
)

func ExitOnBridgeFailedStart(err error) {
	fmt.Printf("An error occurred. %s will now exit. Error: %v\n", shared.PrettyName, err)
	os.Exit(1)
}

func ExitOnDatabaseConnectionFailed(err error) {
	fmt.Printf("Failed to connect to the database. %s will now exit. Error: %v\n", shared.PrettyName, err)
	os.Exit(1)
}

func ExitOnMigrationFailed(err error) {
	fmt.Printf("Database migration failed. %s will now exit. Error: %v\n", shared.PrettyName, err)
	os.Exit(1)
}
