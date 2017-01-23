package goose

import (
	"database/sql"
	"fmt"
)

// SqlDialect abstracts the details of specific SQL dialects
// for goose's few SQL specific statements
type SqlDialect interface {
	createVersionTableSql() string // sql string to create the goose_db_version table
	insertVersionSql() string      // sql string to insert the initial version table row
	dbVersionQuery(db *sql.DB) (*sql.Rows, error)
	versionStatusQuery(db *sql.DB, version int64) (bool, error) //Determines the current IsApplied value of the supplied version
}

var dialect SqlDialect = &PostgresDialect{}

func GetDialect() SqlDialect {
	return dialect
}

func SetDialect(d string) error {
	switch d {
	case "postgres":
		dialect = &PostgresDialect{}
	case "mysql":
		dialect = &MySqlDialect{}
	case "sqlite3":
		dialect = &Sqlite3Dialect{}
	default:
		return fmt.Errorf("%q: unknown dialect", d)
	}

	return nil
}

////////////////////////////
// Postgres
////////////////////////////

type PostgresDialect struct{}

func (pg PostgresDialect) createVersionTableSql() string {
	return `CREATE TABLE goose_db_version (
            	id serial NOT NULL,
                version_id bigint NOT NULL,
                is_applied boolean NOT NULL,
                tstamp timestamp NULL default now(),
                PRIMARY KEY(id)
            );`
}

func (pg PostgresDialect) insertVersionSql() string {
	return "INSERT INTO goose_db_version (version_id, is_applied) VALUES ($1, $2);"
}

func (pg PostgresDialect) dbVersionQuery(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query("SELECT version_id, is_applied from goose_db_version ORDER BY id DESC")
	if err != nil {
		return nil, err
	}

	return rows, err
}

func (pg PostgresDialect) versionStatusQuery(db *sql.DB, version int64) (bool, error) {
	var isApplied bool
	err := db.QueryRow("SELECT is_applied FROM goose_db_version WHERE version_id = $1", version).Scan(&isApplied)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
		return isApplied, nil
	}
}

////////////////////////////
// MySQL
////////////////////////////

type MySqlDialect struct{}

func (m MySqlDialect) createVersionTableSql() string {
	return `CREATE TABLE goose_db_version (
                id serial NOT NULL,
                version_id bigint NOT NULL,
                is_applied boolean NOT NULL,
                tstamp timestamp NULL default now(),
                PRIMARY KEY(id)
            );`
}

func (m MySqlDialect) insertVersionSql() string {
	return "INSERT INTO goose_db_version (version_id, is_applied) VALUES (?, ?);"
}

func (m MySqlDialect) dbVersionQuery(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query("SELECT version_id, is_applied from goose_db_version ORDER BY id DESC")
	if err != nil {
		return nil, err
	}

	return rows, err
}

func (m MySqlDialect) versionStatusQuery(db *sql.DB, version int64) (bool, error) {
	var isApplied bool
	err := db.QueryRow("SELECT is_applied FROM goose_db_version WHERE version_id = ?", version).Scan(&isApplied)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
		return isApplied, nil
	}
}

////////////////////////////
// sqlite3
////////////////////////////

type Sqlite3Dialect struct{}

func (m Sqlite3Dialect) createVersionTableSql() string {
	return `CREATE TABLE goose_db_version (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                version_id INTEGER NOT NULL,
                is_applied INTEGER NOT NULL,
                tstamp TIMESTAMP DEFAULT (datetime('now'))
            );`
}

func (m Sqlite3Dialect) insertVersionSql() string {
	return "INSERT INTO goose_db_version (version_id, is_applied) VALUES (?, ?);"
}

func (m Sqlite3Dialect) dbVersionQuery(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query("SELECT version_id, is_applied from goose_db_version ORDER BY id DESC")
	if err != nil {
		return nil, err
	}

	return rows, err
}

func (m Sqlite3Dialect) versionStatusQuery(db *sql.DB, version int64) (bool, error) {
	var isApplied bool
	err := db.QueryRow("SELECT is_applied FROM goose_db_version WHERE version_id = ?", version).Scan(&isApplied)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
		return isApplied, nil
	}
}
