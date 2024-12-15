package settings

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/mreliasen/scrolls-cli/internal/settings"
	"github.com/mreliasen/scrolls-cli/internal/utils"
)

var library *Library

func New() *Library {
	if library != nil {
		return library
	}

	library = &Library{
		scrolls: map[string]*ScrollMeta{},
	}

	return library
}

type ScrollMeta struct {
	Id       string
	Name     string
	Type     string
	FileName string
}

type Library struct {
	scrolls map[string]*ScrollMeta
	mu      sync.Mutex
}

func (s *Library) getStateFilePath() (string, error) {
	configDir, err := settings.GetConfigDir()
	if err != nil {
		return "", err
	}

	filePath := path.Join(configDir, "library.json")

	if abs, err := filepath.Abs(filePath); err == nil {
		filePath = abs
	}

	return filePath, nil
}

func (s *Library) Load() error {
	if s.scrolls != nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	filePath, err := s.getStateFilePath()
	if err != nil {
		return err
	}

	_, err = os.Stat(filePath)
	if !os.IsNotExist(err) {
		sb, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		data, err := utils.Unmarshal[scrollsState](sb)
		s.state = &data
	}

	fmt.Printf("%+v", s.state)

	return nil
}

func (c *Library) Set(scroll_name, file_name string) {
	if c.state == nil {
		c.Load()

		if c.state == nil {
			panic("failed to update state, state not initialised")
		}
	}

	(*c.state)[scroll_name] = file_name
}

func (c *Library) Get(scroll_name string) string {
	if c.state == nil {
		c.Load()

		if c.state == nil {
			panic("failed to update state, state not initialised")
		}
	}

	return (*c.state)[scroll_name]
}

func (c *Library) PersistChanges() {
	if c.state == nil || !c.changed {
		return
	}

	bytes, err := utils.Marshal(c.state)
	if err != nil {
		fmt.Printf("failed to persist state changes: %s", err.Error())
	}

	filePath, err := c.getStateFilePath()
	if err != nil {
		fmt.Printf("failed to get state file path: %s", err.Error())
	}

	os.WriteFile(filePath, bytes, 0o644)
}
