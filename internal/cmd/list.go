package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/mreliasen/scrolls-cli/internal/flags"
	"github.com/mreliasen/scrolls-cli/internal/scrolls"
	"github.com/mreliasen/scrolls-cli/internal/scrolls/file_handler"
	"github.com/mreliasen/scrolls-cli/internal/tui"
	"github.com/spf13/cobra"
)

// flags for type

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Show a list of all scroll",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		c, err := scrolls.New()
		if err != nil {
			fmt.Printf("failed to initiate scroll-cli")
			return
		}

		st := flags.ScrollType()
		list, err := c.Files.ListScrolls(strings.ToLower(st))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		selection, cancel := tui.NewScrollList(list)
		if cancel {
			return
		}

		action := ""
		form := huh.NewSelect[string]().
			Title(fmt.Sprintf("Selected Scroll: %s", selection.Name)).
			Options(
				huh.NewOption("Edit", "edit"),
				huh.NewOption("Delete", "delete"),
				huh.NewOption("Cancel", "cancel"),
			).
			Value(&action)

		err = form.Run()
		if err != nil {
			log.Fatal(err)
		}

		if action == "cancel" {
			return
		}

		switch action {
		case "edit":
			c.Files.EditScroll(selection.Name)
		case "delete":
			scrollDelete(selection)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	flags.AddScrollTypeFlag(listCmd)
}

func scrollDelete(f *file_handler.FileHandler) {
	confirm := false

	form := huh.NewConfirm().
		Title(fmt.Sprintf("Confirm you want to DELETE the scroll: %s?", f.Name)).
		Affirmative("Yes").
		Negative("No").
		Value(&confirm)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	if !confirm {
		return
	}

	f.Delete()
}
