package scrolls

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/mreliasen/scrolls-cli/internal/file_types"
	"github.com/mreliasen/scrolls-cli/internal/flags"
	"github.com/mreliasen/scrolls-cli/internal/library"
	"github.com/mreliasen/scrolls-cli/internal/settings"
	"github.com/mreliasen/scrolls-cli/internal/tui"
)

func tmpStoragePath() (string, error) {
	p, err := settings.GetConfigDir()
	if err != nil {
		return "", err
	}

	p = path.Join(p, ".tmp")
	err = os.MkdirAll(p, 0o755)
	return p, err
}

type StorageClient client

func (c *StorageClient) Get(name string) (*library.Scroll, error) {
	return c.client.Library.GetByName(name)
}

func (c *StorageClient) Delete(name string) error {
	return c.client.Library.Delete(name)
}

func (c *StorageClient) Rename(src, dist string) error {
	if c.client.Library.Exists(dist) {
		return fmt.Errorf("a scroll already exists with the name %s", dist)
	}

	return c.client.Library.Rename(src, dist)
}

func (c *StorageClient) List() ([]*library.Scroll, error) {
	return nil, nil
}

func (c *StorageClient) New(name string, useTemplate bool, fromFile string) (*library.Scroll, error) {
	path, err := tmpStoragePath()
	if err != nil {
		return nil, err
	}

	if c.client.Library.Exists(name) {
		ok := tui.NewConfirm(fmt.Sprintf("A scroll already exists with the name %s, overwrite?", name))
		if !ok {
			return nil, nil
		}
	}

	templateContent := []byte{}
	if fromFile != "" {
		ffbyte, err := os.ReadFile(fromFile)
		if err != nil {
			return nil, errors.New("failed to read the content of \"%s\"; is the path correct?")
		}

		templateContent = ffbyte
	}

	fType, cancel := tui.NewSelector("")
	if cancel {
		return nil, nil
	}

	tmpId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	ex := file_types.ExecList[fType]
	f := NewFile(path, tmpId.String(), ex.Ext)
	f.Type = fType

	if len(templateContent) == 0 && useTemplate {
		templateContent = []byte(ex.Template)
	}

	if len(templateContent) > 0 {
		f.Write(templateContent)
	}

	f, err = c.editFile(f, nil)
	if err != nil {
		return nil, err
	}

	s, err := c.client.Library.NewScroll(name, fType, f.Body())
	if err == nil {
		fmt.Println("Scroll created!")
	}

	return s, err
}

func (c *StorageClient) NewTempFile(scroll *library.Scroll) (*FileHandler, error) {
	path, err := tmpStoragePath()
	if err != nil {
		return nil, err
	}

	tmpId := uuid.New()
	ex := scroll.Exec()
	f := NewFile(path, tmpId.String(), ex.Ext)
	f.Type = scroll.Type()
	f.Write(scroll.Body())

	return f, nil
}

func (c *StorageClient) EditText(name string) error {
	path, err := tmpStoragePath()
	if err != nil {
		return err
	}

	scroll, err := c.client.Library.GetByName(name)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	ex := file_types.ExecList[scroll.Type()]
	f := NewFile(path, scroll.Id(), ex.Ext)
	f.Type = scroll.Type()
	f.Write(scroll.Body())

	f, err = c.editFile(f, scroll)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	scroll.SetBody(f.Body())
	err = scroll.Save()

	if err == nil {
		fmt.Println("Scroll updated!")
	}

	return err
}

func (c *StorageClient) editFile(f *FileHandler, scroll *library.Scroll) (*FileHandler, error) {
	bin := c.client.Settings.GetEditor()

	cmd := exec.Command(bin, f.Path())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("opening in: %s, waiting to editor to close before proceeding..\n", bin)

	now := time.Now().Add(time.Second * time.Duration(2)).Unix()
	err := cmd.Run()
	if err != nil {
		f.Delete()
		return f, fmt.Errorf("editor error: %s", err.Error())
	}

	end := time.Now().Unix()

	if flags.Debug() {
		fmt.Printf("opening in external editor \"%s\": %t\n", bin, now > end)
	}

	if now > end {
		fmt.Println("When you are done, Press Enter to continue..")
		bufio.NewReader(os.Stdin).ReadString('\n')
	}

	_, err = f.Read()
	if err != nil {
		return f, err
	}

	if scroll != nil {
		if bytes.Equal(f.Body(), scroll.Body()) {
			return f, nil
		}
	}

	f.Delete()
	return f, nil
}
