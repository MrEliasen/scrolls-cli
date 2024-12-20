package scrolls

import "github.com/mreliasen/scrolls-cli/internal/file_types"

type ExecCommand struct {
	Exec     file_types.ExecArgs
	TempFile *FileHandler
}
