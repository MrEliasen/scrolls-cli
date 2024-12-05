package file_handler

type ExecCommand struct {
	Bin      string
	Args     []string
	TempFile *FileHandler
}

type ExecArgs []string

// bin:extension
var ExecFileRequired = map[string]string{
	"go": "go",
}

var ExecList = map[string]ExecArgs{
	"plain-text": {},
	"php": {
		"-r",
	},
	"go": {
		"run",
	},
	"bash": {
		"-c",
	},
	"python": {
		"-c",
	},
	"node": {
		"-e",
	},
	"ruby": {
		"-e",
	},
	"perl": {
		"-e",
	},
	"Rscript": {
		"-e",
	},
	"julia": {
		"-e",
	},
	"cargo": {
		"script",
		"-e",
	},
	"runhashell": {
		"-e",
	},
	"lua": {
		"-e",
	},
	"kotlin": {
		"-script",
	},
	"java": {},
	"powershell": {
		"-command",
	},
	"dotnet": {
		"script",
		"-e",
	},
}
