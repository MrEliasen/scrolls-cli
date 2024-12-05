package cmd

import (
	"fmt"

	"github.com/mreliasen/scrolls-cli/internal/scrolls"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "delete a scroll.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c, err := scrolls.New()
		if err != nil {
			fmt.Printf("failed to initiate scroll-cli")
			return
		}

		c.Files.DeleteScroll(name)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
