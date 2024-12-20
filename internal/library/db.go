package library

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var scrollDBClient *ScrollsDB

func NewConnection(configDir string) (*ScrollsDB, error) {
	if scrollDBClient != nil {
		return scrollDBClient, nil
	}

	dbFile := path.Join(configDir, "scrolls.db")

	if abs, err := filepath.Abs(dbFile); err == nil {
		dbFile = abs
	}

	scrollDBClient = &ScrollsDB{
		DbFile: dbFile,
	}

	err := scrollDBClient.Connect()
	if err != nil {
		return nil, err
	}

	return scrollDBClient, nil
}

type ScrollsDB struct {
	Db     *sql.DB
	DbFile string
	BkFile string
}

func (s *ScrollsDB) Connect() error {
	// check the db file exists
	_, err := os.Stat(s.DbFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		// create a new db file if one does not exist
		if _, err := os.Create(s.DbFile); err != nil {
			return err
		}
	}

	dbConnStr := path.Join("file:", s.DbFile)
	dbClient, err := sql.Open("sqlite3", dbConnStr)
	if err != nil {
		return fmt.Errorf("failed to load scrolls state: %s", err)
	}

	s.Db = dbClient
	return nil
}

func (s *ScrollsDB) Backup() error {
	src, err := os.Stat(scrollDBClient.DbFile)
	if err != nil {
		return err
	}

	// ignore dirs, symlinks etc
	if !src.Mode().IsRegular() {
		return fmt.Errorf("failed to backup scrolls db, the files is not a regular file: %s", scrollDBClient.DbFile)
	}

	scrollDBClient.Db.Close()

	from, err := os.ReadFile(scrollDBClient.DbFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(scrollDBClient.DbFile+".bk", from, 0o644)
	if err != nil {
		return err
	}

	s.BkFile = scrollDBClient.DbFile + ".bk"

	s.Connect()
	return err
}

func (s *ScrollsDB) Restore() error {
	_, err := os.Stat(s.BkFile)
	if err != nil {
		return err
	}

	scrollDBClient.Db.Close()

	err = os.Remove(scrollDBClient.DbFile)
	if err != nil {
		return err
	}

	err = os.Rename(s.BkFile, s.DbFile)
	if err != nil {
		return err
	}

	s.Connect()
	return nil
}

func (s *ScrollsDB) RemoveBackup() error {
	_, err := os.Stat(s.BkFile)
	if err != nil {
		return err
	}

	err = os.Remove(s.BkFile)
	if err != nil {
		return err
	}

	return err
}
