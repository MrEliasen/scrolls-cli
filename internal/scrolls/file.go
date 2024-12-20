package scrolls

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/mreliasen/scrolls-cli/internal/flags"
)

func NewFile(file_path, file_name, file_ext string) *FileHandler {
	name := fmt.Sprintf("%s%s", file_name, file_ext)

	f := &FileHandler{
		path: path.Join(file_path, strings.ToLower(name)),
	}

	return f
}

type FileHandler struct {
	Type string
	data []byte
	path string
}

func (f *FileHandler) Path() string {
	return f.path
}

func (f *FileHandler) Body() []byte {
	return f.data
}

func (f *FileHandler) Delete() error {
	err := os.Remove(f.path)
	if err != nil {
		return fmt.Errorf("failed to delete temporary file\"%s\", %w", f.path, err)
	}

	if flags.Debug() {
		fmt.Printf("deleted temporary file: %s\n", f.path)
	}

	return nil
}

func (f *FileHandler) Write(payload []byte) error {
	f.data = payload
	return os.WriteFile(f.path, payload, 0o644)
}

func (f *FileHandler) Read() ([]byte, error) {
	b, err := os.ReadFile(f.path)
	if err != nil {
		return nil, err
	}

	f.data = b
	return b, nil
}
