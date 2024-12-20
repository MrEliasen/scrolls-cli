package shared

import "database/sql"

func NewMigration(name, upSql, downSql string, up, down func(*sql.DB) error) MigrationInterface {
	return &Migration{name, upSql, downSql, up, down}
}

type Migration struct {
	name    string
	upSql   string
	downSql string
	up      func(*sql.DB) error
	down    func(*sql.DB) error
}

func (m *Migration) Name() string {
	return m.name
}

func (m *Migration) UpSQL() string {
	return m.upSql
}

func (m *Migration) DownSQL() string {
	return m.downSql
}

func (m *Migration) Up(db *sql.DB, migrate func(MigrationInterface) error) error {
	err := migrate(m)
	if err != nil {
		return nil
	}

	if m.up != nil {
		return m.up(db)
	}

	return nil
}

func (m *Migration) Down(db *sql.DB, migrate func(MigrationInterface) error) error {
	err := migrate(m)
	if err != nil {
		return nil
	}

	if m.down != nil {
		return m.down(db)
	}

	return nil
}
