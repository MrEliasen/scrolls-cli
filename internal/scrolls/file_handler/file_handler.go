package file_handler

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mreliasen/scrolls-cli/internal/flags"
)

const headerEndIndicator = "---SCROLL META END---"

func New(file_path string) *FileHandler {
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
		Bin: f.Type,
	}

	if ext, found := ExecFileRequired[f.Type]; found {
		exec.TempFile = f.makeTempFile(ext)
		if exec.TempFile == nil {
			return nil
		}
	}

	for _, a := range ExecList[f.Type] {
		exec.Args = append(exec.Args, a)
	}

	if exec.TempFile != nil {
		exec.Args = append(exec.Args, exec.TempFile.path)
	} else {
		exec.Args = append(exec.Args, f.Body())
	}

	return exec
}

func (f *FileHandler) makeTempFile(ext string) *FileHandler {
	b := f.Body()

	if b == "" {
		return nil
	}

	tmp := New(fmt.Sprintf("%s_%d.%s", f.path, time.Now().UnixMilli(), ext))
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

func (f *FileHandler) WriteHeader() {
	lines, err := f.loadFile()
	if err != nil {
		panic(err)
	}

	payload := []byte(f.generateFileHeader())
	payload = append(payload, []byte("\n")...)
	payload = append(payload, []byte(strings.Join(lines, "\n"))...)

	err = os.WriteFile(f.path, payload, 0o644)
	if err != nil {
		panic(err)
	}

	if flags.Debug() {
		fmt.Printf("wrote scroll %s, with %d bytes\n", f.path, len(payload))
	}
}

func (f *FileHandler) Delete() {
	err := os.Remove(f.path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to delete scroll: %s\n", f.path)
		return
	}

	if flags.Debug() {
		fmt.Printf("deleted cached scrolls: %s\n", f.path)
	}
}

func (f *FileHandler) Save(skipHeader bool) {
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
		panic(err)
	}

	if flags.Debug() {
		fmt.Printf("wrote scroll: %s, with %d bytes\n", f.path, len(payload))
	}
}

func (f *FileHandler) generateFileHeader() string {
	return fmt.Sprintf(
		"%s\nid: %s\nname: %s\ntype: %s\ntags: %s\n\n%s\n%s",
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
		panic(err)
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
