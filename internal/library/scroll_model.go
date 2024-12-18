package library

import (
	"errors"

	"github.com/mreliasen/scrolls-cli/internal/file_types"
)

type Scroll struct {
	uuid      string
	name      string
	file_type string
	body      []byte
	lib       *Library
}

func (s *Scroll) Rename(newName string) error {
	if s.lib.Exists(newName) {
		return errors.New("failed to rename scroll, a scroll already exists with that name")
	}

	s.name = newName
	return s.Save()
}

func (s *Scroll) Save() error {
	return s.lib.Update(s)
}

func (s *Scroll) Delete() error {
	return s.lib.Delete(s.name)
}

func (s *Scroll) Library() *Library {
	return s.lib
}

func (s *Scroll) Id() string {
	return s.uuid
}

func (s *Scroll) Name() string {
	return s.name
}

func (s *Scroll) Type() string {
	return s.file_type
}

func (s *Scroll) Body() []byte {
	return s.body
}

func (s *Scroll) File() string {
	return s.uuid + "." + s.file_type
}

func (s *Scroll) Exec() file_types.ExecArgs {
	return file_types.ExecList[s.Type()]
}

func (s *Scroll) SetType(t string) {
	s.file_type = t
}

func (s *Scroll) SetBody(b []byte) {
	s.body = b
}
