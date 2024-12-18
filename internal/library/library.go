package library

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/mreliasen/scrolls-cli/internal/flags"
	"github.com/mreliasen/scrolls-cli/internal/library/migrations"
	"github.com/mreliasen/scrolls-cli/internal/scrolls/file_handler"
	"github.com/mreliasen/scrolls-cli/internal/settings"
)

var (
	library *Library
	mu      sync.Mutex
)

func LoadLibrary() (*Library, error) {
	mu.Lock()
	defer mu.Unlock()

	if library != nil {
		return library, nil
	}

	cfgDir, err := settings.GetConfigDir()
	if err != nil {
		return nil, err
	}

	db, err := NewConnection(cfgDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load scrolls db: %w", err)
	}

	library = &Library{
		dbClient: db,
		cfgDir:   cfgDir,
	}

	return library, nil
}

type Library struct {
	mu       sync.Mutex
	dbClient *ScrollsDB
	cfgDir   string
}

func (l *Library) ConfigDir() string {
	return l.cfgDir
}

func (l *Library) Migrate() error {
	err := l.dbClient.Backup()
	// backup failed
	if err != nil {
		if flags.Debug() {
			fmt.Printf("%s\n", err.Error())
		}

		// do we have a db connection still?
		fmt.Printf("\n\nFailed to run migrations: failed to backup sqlite db.\n")
		fmt.Printf("You can force the migration without backing up using --skip-backup, however please manually backup your scrolls db first.\n")
		fmt.Printf("DB location: %s\n", l.cfgDir)
		fmt.Printf("Then run: scrolls --skip-backup\n\n")
		return nil
	}

	err = migrations.Migrate(l.dbClient.Db)
	if err != nil {
		l.dbClient.Restore()
		return err
	}

	l.dbClient.RemoveBackup()
	return nil
}

func (l *Library) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.dbClient == nil {
		return
	}

	l.dbClient.Db.Close()
	l.dbClient = nil
}

func (l *Library) NewScroll(scroll_name, file_type string, body []byte) (*Scroll, error) {
	if l.dbClient == nil {
		panic("failed to update scrolls db, db not initialised")
	}

	if l.Exists(scroll_name) {
		return nil, errors.New("a scroll already exists with that name")
	}

	id := uuid.New()

	res, err := l.dbClient.Db.Exec(`
		INSERT INTO
			scrolls (
				uuid,
				name,
				file_type,
				body
			)
		VALUES (?, ?, ?, ?)`,
		id.String(), scroll_name, file_type, []byte(body),
	)
	if err != nil {
		if flags.Debug() {
			fmt.Fprintf(os.Stderr, "%+v", err)
		}
		return nil, err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Scroll{
		uuid:      id.String(),
		name:      scroll_name,
		file_type: file_type,
		body:      body,
		lib:       l,
	}, nil
}

func (l *Library) Exists(name string) bool {
	if l.dbClient == nil {
		panic("failed to update scrolls db, db not initialised")
	}

	res := l.dbClient.Db.QueryRow(`
		SELECT
			uuid
		FROM
			scrolls
		WHERE
			name = ?
	`, name)

	if flags.Debug() && res.Err() != nil {
		fmt.Fprintf(os.Stderr, "%+v", res.Err())
	}

	return res.Err() == nil
}

func (l *Library) Update(scroll *Scroll) error {
	if l.dbClient == nil {
		panic("failed to update scrolls db, db not initialised")
	}

	_, err := l.dbClient.Db.Exec(`
		UPDATE	
			scrolls
		SET
			name = ?,
			file_type = ?,
			body = ?
		WHERE
			uuid = ?
	`, scroll.Name(), scroll.Type(), scroll.Body(), scroll.Id())

	if flags.Debug() && err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
	}

	return err
}

func (l *Library) GetById(id int64) (*Scroll, error) {
	if l.dbClient == nil {
		panic("failed to update scrolls db, db not initialised")
	}

	res := l.dbClient.Db.QueryRow(`
		SELECT
			uuid,
			name,
			file_type,
			body	
		FROM
			scrolls
		WHERE
			id = ?
	`, id)

	scroll := Scroll{
		lib: l,
	}
	err := res.Scan(&scroll.uuid, &scroll.name, &scroll.file_type, &scroll.body)
	if err != nil {
		if flags.Debug() {
			fmt.Fprintf(os.Stderr, "%+v", err)
		}
		return nil, err
	}

	return &scroll, nil
}

func (l *Library) Rename(src, dist string) error {
	if l.dbClient == nil {
		panic("failed to update scrolls db, db not initialised")
	}

	_, err := l.dbClient.Db.Exec(`
		UPDATE
			scrolls
		SET
			name = ?
		WHERE
			name = ?
		LIMIT 1
	`, dist, src)

	if flags.Debug() && err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
	}

	return err
}

func (l *Library) Delete(scroll_name string) error {
	if l.dbClient == nil {
		panic("failed to update scrolls db, db not initialised")
	}

	_, err := l.dbClient.Db.Exec(`
		DELETE FROM	
			scrolls
		WHERE	
			name = ?
	`, scroll_name)

	if flags.Debug() && err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
	}

	return err
}

func (l *Library) GetByName(scroll_name string) (*Scroll, error) {
	if l.dbClient == nil {
		panic("failed to update scrolls db, db not initialised")
	}

	res := l.dbClient.Db.QueryRow(`
		SELECT
			uuid,
			name,
			file_type,
			body
		FROM
			scrolls
		WHERE
			name = ?
	`, scroll_name)

	scroll := Scroll{
		lib: l,
	}
	err := res.Scan(&scroll.uuid, &scroll.name, &scroll.file_type, &scroll.body)
	if err != nil {
		if flags.Debug() {
			fmt.Fprintf(os.Stderr, "%+v", err)
		}
		return nil, err
	}

	return &scroll, nil
}

func (l *Library) GetAllScrollsAutoComplete(name string) ([]*Scroll, error) {
	if l.dbClient == nil {
		return nil, errors.New("scrolls db does not exist")
	}

	res, err := l.dbClient.Db.Query(`
		SELECT
			uuid,
			name,
			file_type,
			body
		FROM
			scrolls
		WHERE
			name LIKE '%' || ?
	`, name)
	if err != nil {
		if flags.Debug() {
			fmt.Fprintf(os.Stderr, "%+v", err)
		}
		return nil, err
	}

	scrolls := []*Scroll{}

	for res.Next() {
		scr := &Scroll{}
		res.Scan(&scr.uuid, &scr.name, &scr.file_type, &scr.body)
		scrolls = append(scrolls, scr)
	}

	return scrolls, nil
}

func (l *Library) GetAllScrollsByType(t string) ([]*Scroll, error) {
	if l.dbClient == nil {
		return nil, errors.New("scrolls db does not exist")
	}

	if t == "all" {
		return l.GetAllScrolls()
	}

	res, err := l.dbClient.Db.Query(`
		SELECT
			uuid,
			name,
			file_type,
			body
		FROM
			scrolls
		WHERE
			file_type = ?
	`, t)
	if err != nil {
		if flags.Debug() {
			fmt.Fprintf(os.Stderr, "%+v", err)
		}
		return nil, err
	}

	scrolls := []*Scroll{}

	for res.Next() {
		scr := &Scroll{}
		res.Scan(&scr.uuid, &scr.name, &scr.file_type, &scr.body)
		scrolls = append(scrolls, scr)
	}

	return scrolls, nil
}

func (l *Library) GetAllScrolls() ([]*Scroll, error) {
	if l.dbClient == nil {
		return nil, errors.New("scrolls db does not exist")
	}

	res, err := l.dbClient.Db.Query(`
		SELECT
			uuid,
			name,
			file_type,
			body
		FROM
			scrolls
	`)
	if err != nil {
		if flags.Debug() {
			fmt.Fprintf(os.Stderr, "%+v", err)
		}
		return nil, err
	}

	scrolls := []*Scroll{}

	for res.Next() {
		scr := &Scroll{}
		res.Scan(&scr.uuid, &scr.name, &scr.file_type, &scr.body)
		scrolls = append(scrolls, scr)
	}

	return scrolls, nil
}

func (l *Library) MigrateScrolls(libPath string) error {
	files, err := os.ReadDir(libPath)
	if err != nil {
		return err
	}

	for _, entry := range files {
		if entry.IsDir() {
			continue
		}

		scroll := file_handler.New(libPath, entry.Name())
		err = scroll.Load()
		if err != nil {
			break
		}

		l.NewScroll(scroll.Name, scroll.Type, []byte(scroll.Body()))
	}

	return err
}
