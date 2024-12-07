package cmd

import (
	"fmt"

	"github.com/mreliasen/scrolls-cli/internal/scrolls"
	"github.com/mreliasen/scrolls-cli/internal/tui"
	"github.com/spf13/cobra"
)

var editorCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a scroll.",
}

var editTextCmd = &cobra.Command{
	Use:               "text <name> ",
	Short:             "Edit a scrolls content.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: ValidScrollName,
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

var editTypeCmd = &cobra.Command{
	Use:               "type <name> ",
	Short:             "Edit a scrolls type.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: ValidScrollName,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c, err := scrolls.New()
		if err != nil {
			fmt.Printf("failed to initiate scroll-cli")
			return
		}

		f, err := c.Files.GetScroll(name)
		if err != nil {
			fmt.Printf("failed to read scroll.")
			return
		}

		f.Type, _ = tui.NewSelector(f.Type)
		f.Save(false)
	},
}

func init() {
	rootCmd.AddCommand(editorCmd)
	editorCmd.AddCommand(editTextCmd, editTypeCmd)
}
