package v0_6

import "github.com/mreliasen/scrolls-cli/internal/library/migrations/shared"

func CreateMigrationTable() shared.MigrationInterface {
	return shared.NewMigration(
		"2024_12_16_1551_create_migrations_table",
		`CREATE TABLE "migrations" (
			"id"	INTEGER,
			"migration"	TEXT NOT NULL,
			PRIMARY KEY("id" AUTOINCREMENT)
		);`,
		`DROP TABLE IF EXISTS "migrations";`,
		nil,
		nil,
	)
}
