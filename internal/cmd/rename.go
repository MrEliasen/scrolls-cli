package cmd

import (
	"fmt"

	"github.com/mreliasen/scrolls-cli/internal/scrolls"
	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:               "rename <name> <new-name>",
	Short:             "Renames a scroll.",
	Aliases:           []string{"mv"},
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: ValidScrollName,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		newName := args[1]

		c, err := scrolls.New()
		if err != nil {
			fmt.Println("failed to initiate scroll-cli")
			return
		}

		err = c.Files.Rename(name, newName)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("Scroll '%s' renamed to '%s'\n", name, newName)
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
