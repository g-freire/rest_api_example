package postgres

import (
	"errors"
	"gym/internal/constants"
	"log"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(databaseURL, migrationsRootFolder, migrationType string, steps int) error {
	m, err := migrate.New(
		migrationsRootFolder,
		databaseURL,
	)
	if err != nil {
		return err
	}

	if migrationType == "up" {
		if steps == 0 {
			err = m.Up()
		} else {
			err = m.Steps(steps)
		}
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		} else {
			log.Print(constants.Green + "MIGRATIONS UP " + strconv.Itoa(steps) + constants.Reset)
		}
	} else if migrationType == "down" {
		if steps == 0 {
			err = m.Down()
		} else {
			err = m.Steps(steps)
		}
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		} else {
			log.Print(constants.Green + "MIGRATIONS DOWN " + strconv.Itoa(steps) + constants.Reset)
		}
	}
	return nil
}
