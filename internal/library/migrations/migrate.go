package migrations

import (
	"database/sql"
	"fmt"

	"github.com/mreliasen/scrolls-cli/internal/flags"
	"github.com/mreliasen/scrolls-cli/internal/library/migrations/shared"
	"github.com/mreliasen/scrolls-cli/internal/library/migrations/v0_6"
)

var migrations = []shared.MigrationInterface{
	v0_6.CreateMigrationTable(),
	v0_6.CreateScrollsTable(),
}

func Migrate(db *sql.DB) error {
	var err error
	failed_migration := ""
	if flags.Debug() {
		fmt.Println("Running Migrations")
	}

	up := migrate(db)
	down := rollback(db)

	for _, m := range migrations {
		if hasMigrationRun(db, m.Name()) {
			if flags.Debug() {
				fmt.Printf("Skipping %s\n", m.Name())
			}
			continue
		}

		if err = m.Up(db, up); err != nil {
			failed_migration = m.Name()
			m.Down(db, down)
			break
		}
	}

	if err != nil {
		err = fmt.Errorf("migration failed: %s\n----\n%w", failed_migration, err)
	}

	if flags.Debug() {
		fmt.Println("Migrations Done")
	}

	return err
}

func hasMigrationRun(db *sql.DB, name string) bool {
	res := db.QueryRow(`
		SELECT
			count(id)	
		FROM
			migrations
		WHERE
			migration = ?
	`, name)

	if res.Err() != nil {
		return false
	}

	var c int64
	res.Scan(&c)

	return c > 0
}

func migrate(db *sql.DB) func(shared.MigrationInterface) error {
	return func(m shared.MigrationInterface) error {
		_, err := db.Exec(m.UpSQL())
		if err != nil {
			return err
		}

		res, err := db.Exec(`
		INSERT INTO
			migrations (migration)
		VALUES
			(?)
	`, m.Name())
		if err != nil {
			return err
		}

		_, err = res.RowsAffected()
		return err
	}
}

func rollback(db *sql.DB) func(shared.MigrationInterface) error {
	return func(m shared.MigrationInterface) error {
		_, err := db.Exec(m.DownSQL())
		if err != nil {
			return err
		}

		_, err = db.Exec(`
		DELETE FROM
			migrations
		WHERE
			migration = ?`,
			m.Name())

		return err
	}
}
