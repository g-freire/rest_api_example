package postgres

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
)

func DoMigrate(migrationsRootFolder, databaseURL string) error {
	m, err := migrate.New(
		migrationsRootFolder,
		databaseURL,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Printf("MIGRATION ERROR: \n", err)
		return err
	}
	log.Printf("MIGRATION WORKS")
	return nil
}
