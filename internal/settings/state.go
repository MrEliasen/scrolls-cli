package settings

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/mreliasen/scrolls-cli/internal/utils"
)

type scrollsState map[string]string

func NewScrollsState() *State {
	return &State{
		state:   &scrollsState{},
		changed: false,
	}
}

type State struct {
	state   *scrollsState
	changed bool
	mu      sync.Mutex
}

func (s *State) getStateFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	filePath := path.Join(configDir, "scrolls.state")

	if abs, err := filepath.Abs(filePath); err == nil {
		filePath = abs
	}

	return filePath, nil
}

func (s *State) Load() error {
	if s.state != nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.state = &scrollsState{}

	filePath, err := GetConfigDir()
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

func (c *State) Set(scroll_name, file_name string) {
	if c.state == nil {
		c.Load()

		if c.state == nil {
			panic("failed to update state, state not initialised")
		}
	}

	(*c.state)[scroll_name] = file_name
}

func (c *State) Get(scroll_name string) string {
	if c.state == nil {
		c.Load()

		if c.state == nil {
			panic("failed to update state, state not initialised")
		}
	}

	return (*c.state)[scroll_name]
}

func (c *State) PersistChanges() {
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
