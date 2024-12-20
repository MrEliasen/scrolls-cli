package cmd

import (
	"github.com/mreliasen/scrolls-cli/internal/library"
	"github.com/spf13/cobra"
)

func ValidScrollName(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	lib, err := library.LoadLibrary()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	names, err := lib.GetAllScrollsAutoComplete(toComplete)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	suggestions := []string{}

	for _, scroll := range names {
		suggestions = append(suggestions, scroll.Name())
	}

	return suggestions, cobra.ShellCompDirectiveNoFileComp
}
