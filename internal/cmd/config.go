package cmd

import (
	"fmt"
	"strings"

	"github.com/mreliasen/scrolls-cli/internal/settings"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage you Scrolls CLI config",
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a config value",
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a config value",
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configSetCmd.AddCommand(setEditorCmd)
	configGetCmd.AddCommand(getEditorCmd)
}

var getEditorCmd = &cobra.Command{
	Use:   "editor",
	Short: "Get the editor which will be used when editing and writing scrolls",
	Run: func(cmd *cobra.Command, args []string) {
		settings, err := settings.LoadSettings()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}

		e := settings.GetEditor()
		if e == "" {
			fmt.Printf("no editor set, using %s as default\n", e)
			return
		}

		s := ""
		if strings.Contains(strings.ToLower(e), "vim") {
			s = ", btw"
		}

		fmt.Printf("scrolls editor: %s%s\n", e, s)
	},
}

var setEditorCmd = &cobra.Command{
	Use:   "editor <name/path>",
	Short: "Set the editor to use when editing and writing scrolls",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{}, cobra.ShellCompDirectiveDefault
	},
	Run: func(cmd *cobra.Command, args []string) {
		settings, err := settings.LoadSettings()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}

		v := strings.Trim(args[0], " ")

		if v == "" {
			fmt.Println("you must specify an editor")
			return
		}

		err = settings.SetEditor(v)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		settings.PersistChanges()

		s := ""
		if strings.Contains(strings.ToLower(v), "vim") {
			s = ", btw"
		}

		fmt.Printf("editor updated to: %s%s\n", v, s)
	},
}
