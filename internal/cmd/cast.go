package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/mreliasen/scrolls-cli/internal/flags"
	"github.com/mreliasen/scrolls-cli/internal/scrolls"
	"github.com/mreliasen/scrolls-cli/internal/scrolls/file_handler"
	"github.com/spf13/cobra"
)

var castCmd = &cobra.Command{
	Use:               "cast <scroll name>",
	Short:             "Run/Execute the scroll.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: ValidScrollName,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c, err := scrolls.New()
		if err != nil {
			fmt.Printf("failed to initiate scroll-cli")
			return
		}

		s, err := c.Files.GetScroll(name)
		if err != nil {
			if flags.Debug() {
				fmt.Printf("%+v\n", err)
			}

			fmt.Printf("failed to retrieve scroll.")
			return
		}

		if s.Type == "plain-text" {
			fmt.Printf("unable to cast a plain text scroll")
			return
		}

		ex := s.GetExec()
		if ex == nil {
			fmt.Printf("no type set for this scroll")
			return
		}

		if !runAsFile(s, ex) {
			err = castInline(s, ex)
			if err == nil {
				return
			}
		}

		tmp := s.MakeTempFile(ex.Exec.Ext)
		if tmp == nil {
			log.Fatalln("Error casting scroll, failed prepare the scroll")
			return
		}

		err = castFile(tmp, ex)
		if err != nil {
			log.Fatalf("Error casting scroll: %s", err.Error())
		}

		tmp.Delete()
	},
}

func runAsFile(s *file_handler.FileHandler, ex *file_handler.ExecCommand) bool {
	if ex.Exec.FileOnly {
		return true
	}

	// some hacky overrides
	switch ex.Exec.Bin {
	case "php":
		for _, l := range s.Lines {
			if strings.Contains(l, "<?") {
				return true
			}
		}
	}

	return false
}

func castInline(s *file_handler.FileHandler, ex *file_handler.ExecCommand) error {
	args := ex.Exec.Args
	args = append(args, s.Body())

	scroll := exec.Command(ex.Exec.Bin, args...)
	scroll.Stdout = os.Stdout
	scroll.Stderr = os.Stderr

	return scroll.Run()
}

func castFile(s *file_handler.FileHandler, ex *file_handler.ExecCommand) error {
	scroll := exec.Command(ex.Exec.Bin, s.Path())
	scroll.Stdout = os.Stdout
	scroll.Stderr = os.Stderr

	return scroll.Run()
}

func init() {
	rootCmd.AddCommand(castCmd)
}
