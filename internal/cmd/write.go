package cmd

import (
	"fmt"

	"github.com/mreliasen/scrolls-cli/internal/flags"
	"github.com/mreliasen/scrolls-cli/internal/scrolls"
	"github.com/spf13/cobra"
)

var writeCmd = &cobra.Command{
	Use:   "write <name>",
	Short: "Create a new scroll, with the given name.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c, err := scrolls.New()
		if err != nil {
			fmt.Println("failed to initiate scroll-cli")
			return
		}

		c.Files.NewScroll(name, flags.Template(), flags.FromFile())
	},
}

func init() {
	rootCmd.AddCommand(writeCmd)
	flags.AddTemplateFlag(writeCmd)
	flags.AddFromFileFlag(writeCmd)
}
