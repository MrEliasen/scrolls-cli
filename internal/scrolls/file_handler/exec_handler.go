package file_handler

type ExecCommand struct {
	Exec     ExecArgs
	TempFile *FileHandler
}

type ExecArgs struct {
	Bin      string
	Args     []string
	Ext      string
	FileOnly bool
}

var ExecList = map[string]ExecArgs{
	"plain-text": {
		Bin:      "",
		Args:     []string{""},
		Ext:      ".txt",
		FileOnly: false,
	},
	"php": {
		Bin:      "php",
		Args:     []string{"-r"},
		Ext:      ".php",
		FileOnly: false,
	},
	"go": {
		Bin:      "go",
		Args:     []string{"run"},
		Ext:      ".go",
		FileOnly: true,
	},
	"bash": {
		Bin:      "bash",
		Args:     []string{"-c"},
		Ext:      ".sh",
		FileOnly: false,
	},
	"python": {
		Bin:      "python",
		Args:     []string{"-c"},
		Ext:      ".py",
		FileOnly: false,
	},
	"node": {
		Bin:      "node",
		Args:     []string{"-e"},
		Ext:      ".js",
		FileOnly: false,
	},
	"ruby": {
		Bin:      "ruby",
		Args:     []string{"-e"},
		Ext:      ".rb",
		FileOnly: false,
	},
	"perl": {
		Bin:      "perl",
		Args:     []string{"-e"},
		Ext:      ".pl",
		FileOnly: false,
	},
	"Rscript": {
		Bin:      "Rscript",
		Args:     []string{"-e"},
		Ext:      ".R",
		FileOnly: false,
	},
	"julia": {
		Bin:      "julia",
		Args:     []string{"-e"},
		Ext:      ".jl",
		FileOnly: false,
	},
	"cargo": {
		Bin:      "cargo",
		Args:     []string{"script", "-e"},
		Ext:      ".rs",
		FileOnly: false,
	},
	"runhashell": {
		Bin:      "runhashell",
		Args:     []string{"-e"},
		Ext:      ".hs",
		FileOnly: false,
	},
	"lua": {
		Bin:      "lua",
		Args:     []string{"-e"},
		Ext:      ".lua",
		FileOnly: false,
	},
	"kotlin": {
		Bin:      "kotlinc",
		Args:     []string{"-script"},
		Ext:      ".kts",
		FileOnly: false,
	},
	"java": {
		Bin:      "java",
		Args:     []string{},
		Ext:      ".java",
		FileOnly: true,
	},
	"powershell": {
		Bin:      "powershell",
		Args:     []string{"-command"},
		Ext:      ".ps1",
		FileOnly: false,
	},
	"dotnet": {
		Bin:      "dotnet",
		Args:     []string{"script", "-e"},
		Ext:      ".csx",
		FileOnly: false,
	},
}
