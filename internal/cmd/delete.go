package cmd

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
	"github.com/mreliasen/scrolls-cli/internal/scrolls"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:               "delete <name>",
	Short:             "Delete a scroll.",
	Aliases:           []string{"rm"},
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: ValidScrollName,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c, err := scrolls.New()
		if err != nil {
			fmt.Printf("failed to initiate scroll-cli")
			return
		}

		confirm := false

		form := huh.NewConfirm().
			Title(fmt.Sprintf("Confirm you want to delete the scroll \"%s\"?", name)).
			Affirmative("Yes").
			Negative("No").
			Value(&confirm)

		err = form.Run()
		if err != nil {
			log.Fatal(err)
		}

		if !confirm {
			return
		}

		c.Files.DeleteScroll(name)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
