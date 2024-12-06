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
		panic(err)
	}

	lib := p.GetLibrary()
	err = os.MkdirAll(lib, 0o755)
	return lib, nil
}

func (c *FileClient) GetScroll(name string) *file_handler.FileHandler {
	path, err := storagePath()
	if err != nil {
		panic(err)
	}

	f := file_handler.New(fmt.Sprintf("%s/%s", path, name))
	if err := f.Load(); err != nil {
		panic(err)
	}

	return f
}

func (c *FileClient) NewScroll(name string, useTemplate bool) {
	path, err := storagePath()
	if err != nil {
		panic(err)
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	fType, cancel := tui.NewSelector("")
	if cancel {
		return
	}
	ex := file_handler.ExecList[fType]

	f := file_handler.New(fmt.Sprintf("%s/%s%s", path, name, ex.Ext))
	f.Id = uuid.String()
	f.Name = name
	f.Type = fType

	if useTemplate {
		os.WriteFile(f.Path(), []byte(ex.Template), 0o644)
	}

	editor := c.client.settings.GetEditor()
	cmd := exec.Command(editor, f.Path())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("opening in: %s, waiting to editor to close before proceeding..\n", c.client.settings.GetEditor())
	err = cmd.Run()
	if err != nil {
		f.Delete()
		log.Fatalf("editor error: %s", err.Error())
		return
	}

	f.WriteHeader()
	os.Rename(f.Path(), fmt.Sprintf("%s/%s", path, name))
}

func (c *FileClient) EditScroll(name string) {
	path, err := storagePath()
	if err != nil {
		panic(err)
	}

	f := file_handler.New(fmt.Sprintf("%s/%s", path, name))
	if !f.Exists() {
		log.Fatalf("no scroll with name \"%s\" found.", name)
		return
	}

	if err := f.Load(); err != nil {
		panic(err)
	}

	tmp_file := f.MakeTempFile(f.GetExec().Exec.Ext)
	createdAt, err := os.Stat(tmp_file.Path())
	if err != nil {
		log.Fatalln("failed to prepare scroll for editing")
		return
	}

	editor := c.client.settings.GetEditor()
	cmd := exec.Command(editor, tmp_file.Path())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatalf("editor error; %s", err.Error())
		return
	}

	updatedAt, err := os.Stat(tmp_file.Path())
	if err != nil {
		log.Fatalln("failed to prepare scroll for editing")
		tmp_file.Delete()
		return
	}

	if createdAt.ModTime().Unix() == updatedAt.ModTime().Unix() {
		log.Println("no changes made to scroll")
		return
	}

	tmp_file.Load()

	// update original file content
	f.Lines = tmp_file.Lines

	// delete tmp
	tmp_file.Delete()

	f.Type, _ = tui.NewSelector(f.Type)

	// write to original
	f.Save(false)
}

func (c *FileClient) PurgeScrolls() error {
	path, err := storagePath()
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

		err = os.Remove(fmt.Sprintf("%s/%s", path, entry.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to delete scroll: %s\n", entry.Name())
		}
	}

	return nil
}

func (c *FileClient) DeleteScroll(name string) {
	path, err := storagePath()
	if err != nil {
		panic(err)
	}

	f := file_handler.New(fmt.Sprintf("%s/%s", path, name))
	if err := f.Load(); err != nil {
		panic(err)
	}

	f.Delete()
}
