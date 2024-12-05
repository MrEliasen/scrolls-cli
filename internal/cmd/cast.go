package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/mreliasen/scrolls-cli/internal/scrolls"
	"github.com/spf13/cobra"
)

var castCmd = &cobra.Command{
	Use:   "cast <name>",
	Short: "cast/execute a scroll.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c, err := scrolls.New()
		if err != nil {
			fmt.Printf("failed to initiate scroll-cli")
			return
		}

		s := c.Files.GetScroll(name)
		if s.Type == "plain-text" {
			fmt.Printf("unable to cast a plain text scroll")
			return
		}

		execParams := s.GetExec()
		if execParams == nil {
			fmt.Printf("no type set for this scroll")
			return
		}

		scroll := exec.Command(execParams.Bin, execParams.Args...)
		scroll.Stdout = os.Stdout
		scroll.Stderr = os.Stderr

		if err := scroll.Run(); err != nil {
			log.Fatalf("Error casting scroll: %v", err)
		}

		if execParams.TempFile != nil {
			execParams.TempFile.Delete()
		}
	},
}

func init() {
	rootCmd.AddCommand(castCmd)
}
