package cmd

import (
	"fmt"

	"github.com/mreliasen/scrolls-cli/internal/scrolls"
	"github.com/spf13/cobra"
)

var editorCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "edit a scroll.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c, err := scrolls.New()
		if err != nil {
			fmt.Printf("failed to initiate scroll-cli")
			return
		}

		c.Files.EditScroll(name)
	},
}

func init() {
	rootCmd.AddCommand(editorCmd)
}
