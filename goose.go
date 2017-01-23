package goose

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"
)

var (
	duplicateCheckOnce sync.Once
	minVersion         = int64(0)
	maxVersion         = int64((1 << 63) - 1)
)

func Run(command string, db *sql.DB, dir string, args ...string) error {
	switch command {
	case "up":
		if err := Up(db, dir); err != nil {
			return err
		}
	case "up-by-one":
		if err := UpByOne(db, dir); err != nil {
			return err
		}
	case "up-from":
		if len(args) == 0 {
			return fmt.Errorf("up-from must be of form: goose [OPTIONS] DRIVER DBSTRING up-from VERSION")
		}

		v := args[0]
		if len(args[0]) > 14 {
			return fmt.Errorf("VERSION should be less/equal to 14 in length and be in the form YYYYMMDDHHMMSS")
		}

		//Ensure proper length and parse
		v += "00000000000000"
		v = v[0:14]

		version, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse up-from version.  Ensure format is YYYYMMDDHHMMSS")
		}

		if err := UpFrom(db, dir, version); err != nil {
			return err
		}
	case "create":
		if len(args) == 0 {
			return fmt.Errorf("create must be of form: goose [OPTIONS] DRIVER DBSTRING create NAME [go|sql]")
		}

		migrationType := "go"
		if len(args) == 2 {
			migrationType = args[1]
		}
		if err := Create(db, dir, args[0], migrationType); err != nil {
			return err
		}
	case "down":
		if err := Down(db, dir); err != nil {
			return err
		}
	case "redo":
		if err := Redo(db, dir); err != nil {
			return err
		}
	case "status":
		if err := Status(db, dir); err != nil {
			return err
		}
	case "version":
		if err := Version(db, dir); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%q: no such command", command)
	}
	return nil
}
