package file_handler

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mreliasen/scrolls-cli/internal/flags"
	"github.com/mreliasen/scrolls-cli/internal/settings"
)

const headerEndIndicator = "---SCROLL META END---"

func New(file_path, file_name string) *FileHandler {
	f := &FileHandler{
		Id:    "",
		Name:  "",
		Type:  "",
		Tags:  []string{},
		Lines: []string{},
		path:  path.Join(file_path, strings.ToLower(file_name)),
	}

	return f
}

func NewFromFile(file_path string) *FileHandler {
	f := &FileHandler{
		Id:    "",
		Name:  "",
		Type:  "",
		Tags:  []string{},
		Lines: []string{},
		path:  file_path,
	}

	return f
}

type FileHandler struct {
	Id    string
	Name  string
	Type  string
	Tags  []string
	Lines []string
	path  string
}

func (f *FileHandler) GetExec() *ExecCommand {
	if f.Type == "" {
		return nil
	}

	exec := &ExecCommand{
		Exec: ExecList[f.Type],
	}

	return exec
}

func (f *FileHandler) MakeTempFile(ext string) *FileHandler {
	b := f.Body()

	if b == "" {
		return nil
	}

	configDir, err := settings.GetConfigDir()
	if err != nil {
		return nil
	}

	tmp_path := path.Join(configDir, "/tmp")
	err = os.MkdirAll(tmp_path, 0o755)
	if err != nil {
		return nil
	}

	tmp_path = path.Join(tmp_path, f.Name)
	tmp := NewFromFile(fmt.Sprintf("%s_%d%s", tmp_path, time.Now().UnixMilli(), ext))
	tmp.Lines = f.Lines
	tmp.Save(true)
	return tmp
}

func (f *FileHandler) Path() string {
	return f.path
}

func (f *FileHandler) Body() string {
	return strings.Join(f.Lines, "\n")
}

func (f *FileHandler) Exists() bool {
	_, err := os.Stat(f.path)
	return err == nil
}

func (f *FileHandler) WriteHeader() error {
	lines, err := f.loadFile()
	if err != nil {
		return fmt.Errorf("failed to load scrolls, %w", err)
	}

	payload := []byte(f.generateFileHeader())
	payload = append(payload, []byte("\n")...)
	payload = append(payload, []byte(strings.Join(lines, "\n"))...)

	err = os.WriteFile(f.path, payload, 0o644)
	if err != nil {
		return fmt.Errorf("failed to save scroll, %w", err)
	}

	if flags.Debug() {
		fmt.Printf("wrote scroll %s, with %d bytes\n", f.path, len(payload))
	}

	return nil
}

func (f *FileHandler) Delete() error {
	err := os.Remove(f.path)
	if err != nil {
		return fmt.Errorf("failed to delete scroll \"%s\", %w\n", f.path, err)
	}

	if flags.Debug() {
		fmt.Printf("deleted scroll: %s\n", f.path)
	}

	return nil
}

func (f *FileHandler) Rename(new_name string) error {
	scrollsDir := path.Dir(f.path)
	f.Name = new_name
	newPath := path.Join(scrollsDir, new_name)

	err := os.Rename(f.path, newPath)
	if err != nil {
		return err
	}

	f.path = newPath
	f.Save(false)
	return nil
}

func (f *FileHandler) Save(skipHeader bool) error {
	header := ""
	if !skipHeader {
		header = f.generateFileHeader() + "\n"
	}

	content := fmt.Sprintf("%s%s",
		header,
		strings.Join(f.Lines, "\n"),
	)

	payload := []byte(content)

	err := os.WriteFile(f.path, payload, 0o644)
	if err != nil {
		return fmt.Errorf("failed to save scroll: %w", err)
	}

	if flags.Debug() {
		fmt.Printf("wrote scroll: %s, with %d bytes\n", f.path, len(payload))
	}

	return nil
}

func (f *FileHandler) generateFileHeader() string {
	return fmt.Sprintf(
		"%s\nid: %s\nname: %s\ntype: %s\ntags: %s\n\n%s\n%s\n",
		"---",
		f.Id,
		f.Name,
		f.Type,
		strings.Join(f.Tags, ", "),
		"Do not edit the this meta data.",
		headerEndIndicator,
	)
}

func (f *FileHandler) loadFile() (lines []string, error error) {
	file, err := os.Open(f.path)
	if err != nil {
		return lines, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for {
		hasLine := scanner.Scan()
		if !hasLine {
			break
		}

		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func (f *FileHandler) Load() error {
	lines, err := f.loadFile()
	if err != nil {
		return fmt.Errorf("failed to load scroll, %w", err)
	}

	f.parse(lines)
	return nil
}

func (f *FileHandler) parse(lines []string) {
	f.Lines = lines

	for i, l := range lines {
		// after the header close tag, we strip the header data from the scroll body
		if l == headerEndIndicator {
			f.Lines = lines[i+1:]
			break
		}

		part := strings.Split(l, ":")
		if len(part) < 2 {
			continue
		}

		key := strings.Trim(part[0], " ")
		val := strings.Trim(part[1], " ")

		switch key {
		case "id":
			f.Id = val
		case "name":
			f.Name = val
		case "type":
			if _, found := ExecList[val]; found {
				f.Type = val
			}
		case "tags":
			tags := strings.Split(val, ",")
			for idx, t := range tags {
				tags[idx] = strings.ToLower(strings.Trim(t, " "))
			}

			f.Tags = tags
		}
	}
}
