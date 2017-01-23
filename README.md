# goose

Goose is a database migration tool. Manage your database's evolution by creating incremental SQL files or Go functions.


### Goals of this fork

This is a fork of https://github.com/pressly/goose with the following change:
Added **up-from** command to execute migrations from a certain point in time.  This is especially useful in a multi-developer 
context where each has a dev database and creates independently dated migrations that are merged to a common repo.

As an example:

Chronologically:
**Migration \#1** Jan 12, 2016 - Bob creates a products table  : 20160112xxxxxx_CreateProductTable.sql
**Migration \#2** Jan 13, 2016 - Mike creates a chat table  : 20160113xxxxxx_CreateChatTable.sql
**Migration \#3** Jan 14, 2016 - Bob creates a orders table which refers to products :  20160114xxxxxx_CreateOrdersTable.sql

If Bob and Mike each subsequently combine their migrations and attempt a GOOSE UP, neither will have the complete database.  In fact, Mike will fail on 
**Migration \#3** because GOOSE UP skipped **Migration \#1**  

The new up-from command allows these migrations to occur.  Mike can simply issue **goose up-from 20160112** (abbrev.) and Goose will simply apply *what hasn't already been applied.*


# Install

    $ go get -u github.com/MatthewLea/goose/cmd/goose

This will install the `goose` binary to your `$GOPATH/bin` directory.

# Usage

```
Usage: goose [OPTIONS] DRIVER DBSTRING COMMAND

Examples:
    goose postgres "user=postgres dbname=postgres sslmode=disable" up
    goose mysql "user:password@/dbname" down
    goose sqlite3 ./foo.db status

Options:
  -dir string
    	directory with migration files (default ".")

Commands:
    up         Migrate the DB to the most recent version available
    up-from    Execute migrations starting from the specified version
    down       Roll back the version by 1
    redo       Re-run the latest migration
    status     Dump the migration status for the current DB
    dbversion  Print the current version of the database
    create     Creates a blank migration template
```
## create

Create a new Go migration.

    $ goose create AddSomeColumns
    $ goose: created db/migrations/20130106093224_AddSomeColumns.go

Edit the newly created script to define the behavior of your migration.

You can also create an SQL migration:

    $ goose create AddSomeColumns sql
    $ goose: created db/migrations/20130106093224_AddSomeColumns.sql

## up

Apply all available migrations.

    $ goose up
    $ goose: migrating db environment 'development', current version: 0, target: 3
    $ OK    001_basics.sql
    $ OK    002_next.sql
    $ OK    003_and_again.go

## up-from

Apply any *unapplied* migrations starting from the specified version (or date/time)

## down

Roll back a single migration from the current version.

    $ goose down
    $ goose: migrating db environment 'development', current version: 3, target: 2
    $ OK    003_and_again.go

## redo

Roll back the most recently applied migration, then run it again.

    $ goose redo
    $ goose: migrating db environment 'development', current version: 3, target: 2
    $ OK    003_and_again.go
    $ goose: migrating db environment 'development', current version: 2, target: 3
    $ OK    003_and_again.go

## status

Print the status of all migrations:

    $ goose status
    $ goose: status for environment 'development'
    $   Applied At                  Migration
    $   =======================================
    $   Sun Jan  6 11:25:03 2013 -- 001_basics.sql
    $   Sun Jan  6 11:25:03 2013 -- 002_next.sql
    $   Pending                  -- 003_and_again.go

## dbversion

Print the current version of the database:

    $ goose dbversion
    $ goose: dbversion 002

# Migrations

goose supports migrations written in SQL or in Go.

## SQL Migrations

A sample SQL migration looks like:

```sql
-- +goose Up
CREATE TABLE post (
    id int NOT NULL,
    title text,
    body text,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE post;
```

Notice the annotations in the comments. Any statements following `-- +goose Up` will be executed as part of a forward migration, and any statements following `-- +goose Down` will be executed as part of a rollback.

By default, SQL statements are delimited by semicolons - in fact, query statements must end with a semicolon to be properly recognized by goose.

More complex statements (PL/pgSQL) that have semicolons within them must be annotated with `-- +goose StatementBegin` and `-- +goose StatementEnd` to be properly recognized. For example:

```sql
-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION histories_partition_creation( DATE, DATE )
returns void AS $$
DECLARE
  create_query text;
BEGIN
  FOR create_query IN SELECT
      'CREATE TABLE IF NOT EXISTS histories_'
      || TO_CHAR( d, 'YYYY_MM' )
      || ' ( CHECK( created_at >= timestamp '''
      || TO_CHAR( d, 'YYYY-MM-DD 00:00:00' )
      || ''' AND created_at < timestamp '''
      || TO_CHAR( d + INTERVAL '1 month', 'YYYY-MM-DD 00:00:00' )
      || ''' ) ) inherits ( histories );'
    FROM generate_series( $1, $2, '1 month' ) AS d
  LOOP
    EXECUTE create_query;
  END LOOP;  -- LOOP END
END;         -- FUNCTION END
$$
language plpgsql;
-- +goose StatementEnd
```

## Go Migrations

Import `github.com/pressly/goose` from your own project (see [example](./example/migrations-go/cmd/main.go)), register migration functions and run goose command (ie. `goose.Up(db *sql.DB, dir string)`).

A [sample Go migration 00002_users_add_email.go file](./example/migrations-go/00002_users_add_email.go) looks like:

```go
package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up, Down)
}

func Up(tx *sql.Tx) error {
	_, err := tx.Query("ALTER TABLE users ADD COLUMN email text DEFAULT '' NOT NULL;")
	if err != nil {
		return err
	}
	return nil
}

func Down(tx *sql.Tx) error {
	_, err := tx.Query("ALTER TABLE users DROP COLUMN email;")
	if err != nil {
		return err
	}
	return nil
}
```

## License

Licensed under [MIT License](./LICENSE)

[GoDoc]: https://godoc.org/github.com/pressly/goose
[GoDoc Widget]: https://godoc.org/github.com/pressly/goose?status.svg
[Travis]: https://travis-ci.org/pressly/goose
[Travis Widget]: https://travis-ci.org/pressly/goose.svg?branch=master
