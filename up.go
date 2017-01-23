package goose

import (
	"database/sql"
	"fmt"
)

func Up(db *sql.DB, dir string) error {
	migrations, err := collectMigrations(dir, minVersion, maxVersion)
	if err != nil {
		return err
	}

	for {
		current, err := GetDBVersion(db)
		if err != nil {
			return err
		}

		next, err := migrations.Next(current)
		if err != nil {
			if err == ErrNoNextVersion {
				fmt.Printf("goose: no migrations to run. current version: %d\n", current)
				return nil
			}
			return err
		}

		if err = next.Up(db); err != nil {
			return err
		}
	}

	return nil
}

func UpByOne(db *sql.DB, dir string) error {
	migrations, err := collectMigrations(dir, minVersion, maxVersion)
	if err != nil {
		return err
	}

	currentVersion, err := GetDBVersion(db)
	if err != nil {
		return err
	}

	next, err := migrations.Next(currentVersion)
	if err != nil {
		if err == ErrNoNextVersion {
			fmt.Printf("goose: no migrations to run. current version: %d\n", currentVersion)
		}
		return err
	}

	if err = next.Up(db); err != nil {
		return err
	}

	return nil
}

func UpFrom(db *sql.DB, dir string, startVersion int64) error {
	migrations, err := collectMigrations(dir, minVersion, maxVersion)
	if err != nil {
		return err
	}

	//find the migration having the lowest version at or above the specified startVersion
	found := false
	for i, v := range migrations {
		if v.Version >= startVersion {
			migrations = migrations[i:len(migrations)]
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("goose: no migrations to run.")
		return nil
	}

	//Iterate through all migrations
	for _, next := range migrations {

		//Only execute if the migration has not already been applied
		isApplied, err := next.GetVersionStatus(db)
		if err != nil {
			fmt.Printf("goose up-from error:%v", err)
			return err
		}

		if isApplied == false {
			if err := next.Up(db); err != nil {
				fmt.Printf("goose up-from error:%v", err)
				return err
			}
		}
	}

	return nil
}
