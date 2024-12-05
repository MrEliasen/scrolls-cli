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

func (c *FileClient) NewScroll(name string) {
	path, err := storagePath()
	if err != nil {
		panic(err)
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	f := file_handler.New(fmt.Sprintf("%s/%s", path, name))
	f.Id = uuid.String()
	f.Name = name

	editor := c.client.settings.GetEditor()
	cmd := exec.Command(editor, f.Path())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatalf("editor error: %s", err.Error())
		return
	}

	f.Type = tui.NewSelector(f.Type)
	f.WriteHeader()
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

	tmp_file := file_handler.New(fmt.Sprintf("%s.scroll_tmp", path))
	tmp_file.Lines = f.Lines
	tmp_file.Save(true)

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

	tmp_file.Load()

	log.Printf("%s", tmp_file.Body())

	// update original file content
	f.Lines = tmp_file.Lines

	// delete tmp
	tmp_file.Delete()

	f.Type = tui.NewSelector(f.Type)

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
