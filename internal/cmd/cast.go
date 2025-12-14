package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/mreliasen/scrolls-cli/internal/flags"
	"github.com/mreliasen/scrolls-cli/internal/library"
	"github.com/mreliasen/scrolls-cli/internal/scrolls"
	"github.com/spf13/cobra"
)

var castCmd = &cobra.Command{
	Use:               "cast <scroll name>",
	Short:             "Run/Execute the scroll.",
	Args:              cobra.ExactArgs(1),
	Aliases:           []string{"run", "exec"},
	ValidArgsFunction: ValidScrollName,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c, err := scrolls.New()
		if err != nil {
			fmt.Println("failed to initiate scroll-cli")
			return
		}

		s, err := c.Storage.Get(name)
		if err != nil {
			if flags.Debug() {
				fmt.Printf("%+v\n", err)
			}

			fmt.Println("failed to retrieve scroll.")
			return
		}

		if s.Type() == "plain-text" {
			fmt.Println("unable to cast a plain text scroll")
			return
		}

		ex := s.Exec()
		if ex.Bin == "" {
			fmt.Println("no type set for this scroll")
			return
		}

		if !runAsFile(s) {
			err = castInline(s)
			if err == nil {
				return
			}
		}

		tmp, err := c.Storage.NewTempFile(s)
		if err != nil {
			log.Fatalln("Error casting scroll, failed prepare the scroll")
			return
		}

		err = castFile(s, tmp.Path())
		if err != nil {
			log.Fatalf("Error casting scroll: %s", err.Error())
		}

		tmp.Delete()
	},
}

func runAsFile(s *library.Scroll) bool {
	ex := s.Exec()
	if ex.FileOnly {
		return true
	}

	// some hacky overrides
	switch ex.Bin {
	case "php":
		if bytes.Contains(s.Body(), []byte("<?")) {
			return true
		}
	}

	return false
}

func castInline(s *library.Scroll) error {
	ex := s.Exec()
	args := append(ex.Args, string(s.Body()))

	scroll := exec.Command(ex.Bin, args...)
	scroll.Stdout = os.Stdout
	scroll.Stderr = os.Stderr

	return scroll.Run()
}

func castFile(s *library.Scroll, path string) error {
	ex := s.Exec()
	args := []string{}

	if len(ex.Args) > 0 {
		args = append(args, ex.Args...)
	}

	if ex.AlwaysUseArgs {
		args = append(ex.Args, path)
	}

	scroll := exec.Command(ex.Bin, args...)
	scroll.Stdout = os.Stdout
	scroll.Stderr = os.Stderr

	return scroll.Run()
}

func init() {
	rootCmd.AddCommand(castCmd)
}
