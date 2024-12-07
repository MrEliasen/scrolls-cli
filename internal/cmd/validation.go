package cmd

import (
	"fmt"
	"os"

	"github.com/mreliasen/scrolls-cli/internal/settings"
	"github.com/spf13/cobra"
)

func ValidScrollName(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	settings, err := settings.LoadSettings()
	if err != nil {
		fmt.Println("unable to read current library path")
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// List files in the current directory
	files, err := os.ReadDir(settings.GetLibrary())
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var suggestions []string
	for _, file := range files {
		if !file.IsDir() && (toComplete == "" || matchesPrefix(file.Name(), toComplete)) {
			suggestions = append(suggestions, file.Name())
		}
	}

	return suggestions, cobra.ShellCompDirectiveNoFileComp
}

func matchesPrefix(name, prefix string) bool {
	return len(name) >= len(prefix) && name[:len(prefix)] == prefix
}
