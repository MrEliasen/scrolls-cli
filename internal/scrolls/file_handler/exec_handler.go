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
	Template string
}

var ExecList = map[string]ExecArgs{
	"plain-text": {
		Bin:      "",
		Args:     []string{""},
		Ext:      ".txt",
		FileOnly: false,
		Template: "",
	},
	"php": {
		Bin:      "php",
		Args:     []string{"-r"},
		Ext:      ".php",
		FileOnly: false,
		Template: `<?php
function main() {
    echo "Hello, World!\n";
}

main();`,
	},
	"go": {
		Bin:      "go",
		Args:     []string{"run"},
		Ext:      ".go",
		FileOnly: true,
		Template: `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`,
	},
	"bash": {
		Bin:      "bash",
		Args:     []string{"-c"},
		Ext:      ".sh",
		FileOnly: false,
		Template: `#!/bin/bash

main() {
  echo "Hello, World!"
}

main`,
	},
	"python": {
		Bin:      "python",
		Args:     []string{"-c"},
		Ext:      ".py",
		FileOnly: false,
		Template: `def main():
print("Hello, World!")

if __name__ == "__main__":
    main()
`,
	},
	"node": {
		Bin:      "node",
		Args:     []string{"-e"},
		Ext:      ".js",
		FileOnly: false,
		Template: `function main() {
	console.log("Hello, World!);
}

main();`,
	},
	"ruby": {
		Bin:      "ruby",
		Args:     []string{"-e"},
		Ext:      ".rb",
		FileOnly: false,
		Template: `def main
  puts "Hello, World!"
end

main`,
	},
	"perl": {
		Bin:      "perl",
		Args:     []string{"-e"},
		Ext:      ".pl",
		FileOnly: false,
		Template: `#!/usr/bin/perl

use strict;
use warnings;

sub main {
    print "Hello, World!\n";
}

main();`,
	},
	"R": {
		Bin:      "Rscript",
		Args:     []string{"-e"},
		Ext:      ".R",
		FileOnly: false,
		Template: `main <- function() {
    print("Hello, World!")
}

main()`,
	},
	"julia": {
		Bin:      "julia",
		Args:     []string{"-e"},
		Ext:      ".jl",
		FileOnly: false,
		Template: `function main()
    println("Hello, World!")
end

main()`,
	},
	"rust": {
		Bin:      "cargo",
		Args:     []string{"script", "-e"},
		Ext:      ".rs",
		FileOnly: false,
		Template: `fn main() {
    println!("Hello, World!");
}`,
	},
	"hashell": {
		Bin:      "runhashell",
		Args:     []string{"-e"},
		Ext:      ".hs",
		FileOnly: false,
		Template: `main :: IO ()
main = putStrLn "Hello, World!"`,
	},
	"lua": {
		Bin:      "lua",
		Args:     []string{"-e"},
		Ext:      ".lua",
		FileOnly: false,
		Template: `local function main()
    print("Hello, World!")
end

main()`,
	},
	"kotlin": {
		Bin:      "kotlinc",
		Args:     []string{"-script"},
		Ext:      ".kts",
		FileOnly: false,
		Template: `fun main() {
    println("Hello, World!")
}

main()`,
	},
	"java": {
		Bin:      "java",
		Args:     []string{},
		Ext:      ".java",
		FileOnly: true,
		Template: `public class HelloWorld {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
    }
}
`,
	},
	"powershell": {
		Bin:      "powershell",
		Args:     []string{"-command"},
		Ext:      ".ps1",
		FileOnly: false,
		Template: `function Main {
    Write-Output "Hello, World!"
}

Main`,
	},
	"dotnet": {
		Bin:      "dotnet",
		Args:     []string{"script", "-e"},
		Ext:      ".csx",
		FileOnly: false,
		Template: `using System;

void Main() {
    Console.WriteLine("Hello, World!");
}

Main();`,
	},
}
