package cmd

type ExecCommand []string

var ExecFileRequired = []string{
	"go",
}

var ExecList = map[string]ExecCommand{
	"php": {
		"-r",
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
