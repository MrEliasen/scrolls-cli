package scrolls

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/google/uuid"
	"github.com/mreliasen/scrolls-cli/internal/scrolls/file_handler"
	"github.com/mreliasen/scrolls-cli/internal/settings"
	"github.com/mreliasen/scrolls-cli/internal/tui"
)

type FileClient client

func storagePath() (string, error) {
	p, err := settings.LoadSettings()
	if err != nil {
		return "", err
	}

	lib := p.GetLibrary()
	err = os.MkdirAll(lib, 0o755)
	return lib, nil
}

func (c *FileClient) GetScroll(name string) (*file_handler.FileHandler, error) {
	path, err := storagePath()
	if err != nil {
		return nil, err
	}

	f := file_handler.New(fmt.Sprintf("%s/%s", path, name))
	if err := f.Load(); err != nil {
		return nil, err
	}

	return f, nil
}

func (c *FileClient) NewScroll(name string, useTemplate bool) error {
	path, err := storagePath()
	if err != nil {
		return err
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	fType, cancel := tui.NewSelector("")
	if cancel {
		return nil
	}
	ex := file_handler.ExecList[fType]

	f := file_handler.New(fmt.Sprintf("%s/%s%s", path, name, ex.Ext))
	f.Id = uuid.String()
	f.Name = name
	f.Type = fType

	if useTemplate {
		os.WriteFile(f.Path(), []byte(ex.Template), 0o644)
	}

	editor := c.client.Settings.GetEditor()
	cmd := exec.Command(editor, f.Path())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("opening in: %s, waiting to editor to close before proceeding..\n", c.client.Settings.GetEditor())
	err = cmd.Run()
	if err != nil {
		f.Delete()
		return fmt.Errorf("editor error: %s", err.Error())
	}

	f.WriteHeader()
	os.Rename(f.Path(), fmt.Sprintf("%s/%s", path, name))

	return nil
}

func (c *FileClient) EditScroll(name string) error {
	path, err := storagePath()
	if err != nil {
		return err
	}

	f := file_handler.New(fmt.Sprintf("%s/%s", path, name))
	if !f.Exists() {
		return fmt.Errorf("no scroll with name \"%s\" found.", name)
	}

	if err := f.Load(); err != nil {
		return err
	}

	tmp_file := f.MakeTempFile(f.GetExec().Exec.Ext)
	createdAt, err := os.Stat(tmp_file.Path())
	if err != nil {
		return fmt.Errorf("failed to prepare scroll for editing")
	}

	editor := c.client.Settings.GetEditor()
	cmd := exec.Command(editor, tmp_file.Path())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("editor error; %s", err.Error())
	}

	updatedAt, err := os.Stat(tmp_file.Path())
	if err != nil {
		tmp_file.Delete()
		return fmt.Errorf("failed to prepare scroll for editing, %w", err)
	}

	if createdAt.ModTime().Unix() == updatedAt.ModTime().Unix() {
		return nil
	}

	tmp_file.Load()

	// update original file content
	f.Lines = tmp_file.Lines

	// delete tmp
	tmp_file.Delete()

	// write to original
	f.Save(false)

	return nil
}

func (c *FileClient) PurgeScrolls() error {
	path, err := storagePath()
	if err != nil {
		return err
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range files {
		if entry.IsDir() {
			continue
		}

		err = os.Remove(fmt.Sprintf("%s/%s", path, entry.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to delete scroll: %s\n", entry.Name())
		}
	}

	return nil
}

func (c *FileClient) DeleteScroll(name string) error {
	path, err := storagePath()
	if err != nil {
		return err
	}

	f := file_handler.New(fmt.Sprintf("%s/%s", path, name))
	if !f.Exists() {
		return fmt.Errorf("the scroll \"%s\" does not exist or is inaccessible.\n", name)
	}

	f.Delete()
	return nil
}

func (c *FileClient) ListScrolls() error {
	/* path, err := storagePath()
	if err != nil {
		panic(err)
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range files {
		if entry.IsDir() {
			continue
		}

		f := file_handler.New(fmt.Sprintf("%s/%s", path, entry.Name()))
	} */

	return nil
}
