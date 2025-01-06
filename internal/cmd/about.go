package cmd

import (
	"fmt"

	"github.com/mreliasen/scrolls-cli/internal/tui"
	"github.com/spf13/cobra"
)

var aboutCmd = &cobra.Command{
	Use:   "about",
	Short: "Display information about Scrolls.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\nScrolls is a CLI tool for making, managing and using snippets/scripts in your terminal. Make a snippet in one of the many supported languages and execute it whenever you need it or simply echo it to stdout.\n\n")
		fmt.Printf("%s:  https://github.com/MrEliasen/scrolls-cli\n", tui.SuccessStyle.Render("Github"))
		fmt.Printf("%s: https://www.oogabooga.dev\n", tui.SuccessStyle.Render("Made By"))
	},
}

func init() {
	rootCmd.AddCommand(aboutCmd)
}
