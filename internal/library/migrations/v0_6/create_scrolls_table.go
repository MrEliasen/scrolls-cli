package v0_6

import (
	"github.com/mreliasen/scrolls-cli/internal/library/migrations/shared"
)

func CreateScrollsTable() shared.MigrationInterface {
	return shared.NewMigration(
		"2024_12_18_1155_create_scrolls_table",
		`CREATE TABLE "scrolls" (
			"id"	INTEGER,
			"uuid"	TEXT UNIQUE,
			"name"	TEXT UNIQUE,
			"file_type"	TEXT,
			"body"	BLOB NOT NULL,
			PRIMARY KEY("id" AUTOINCREMENT)
		);`,
		`DROP TABLE IF EXISTS "scrolls";`,
		nil,
		nil,
	)
}
