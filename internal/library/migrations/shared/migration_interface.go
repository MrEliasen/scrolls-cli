package shared

import "database/sql"

type MigrationInterface interface {
	Name() string
	UpSQL() string
	DownSQL() string
	Up(*sql.DB, func(MigrationInterface) error) error
	Down(*sql.DB, func(MigrationInterface) error) error
}
