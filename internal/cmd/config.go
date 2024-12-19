package cmd

import (
	"fmt"
	"os"
	"path"
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
	configSetCmd.AddCommand(setLibraryCmd)
	configGetCmd.AddCommand(getLibraryCmd)
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
	Short: "Set the editor to use when editing scrolls",
	Long: `Set the editor to use when editing scrolls. This includes terminal and external editors.
Example (vim):	scrolls config set editor vim
Example (zed):	scrolls config set editor zed
Example (path):	scrolls config set editor "/abs/path/to/editor"
	`,
	Args: cobra.ExactArgs(1),
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

var getLibraryCmd = &cobra.Command{
	Use:   "library",
	Short: "Get the path to where your scrolls DB is stored.",
	Run: func(cmd *cobra.Command, args []string) {
		settings, err := settings.LoadSettings()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}

		e := settings.GetLibrary()
		if e == "" {
			fmt.Printf("no library set, using default location %s\n", e)
			return
		}

		fmt.Printf("%s\n", path.Join(e, "scrolls.db"))
	},
}

var setLibraryCmd = &cobra.Command{
	Use:   "library <path>",
	Short: "Set the path to where you want to store you scrolls DB. The existing scrolls DB will be automatically moved to the new location",
	Args:  cobra.MaximumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{}, cobra.ShellCompDirectiveDefault
	},
	Run: func(cmd *cobra.Command, args []string) {
		settings, err := settings.LoadSettings()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}

		v := ""
		if len(args) > 0 {
			v = strings.Trim(args[0], " ")
		}

		src := settings.GetLibrary()
		settings.SetLibrary(v)
		loc := settings.GetLibrary()

		if src == loc {
			fmt.Println("your scrolls location is already set to this location.")
			return
		}

		if v != "" {
			stat, err := os.Stat(v)
			if err != nil {
				if !os.IsNotExist(err) {
					fmt.Println("failed to set new library location.")
					fmt.Println(err.Error())
					return
				}

				err = os.MkdirAll(v, 0o755)
				if err != nil {
					fmt.Println("failed to create the new library folder.")
					return
				}
			} else {
				if !stat.IsDir() {
					fmt.Println("failed to set new library location. The specified location is not a directory.")
					return
				}
			}
		}

		_, err = os.Stat(path.Join(loc, "scrolls.db"))
		if err == nil {
			fmt.Println("failed to set new library location. There is already a file called \"scrolls.db\" in the new location.")
			return
		}

		err = os.Rename(
			path.Join(src, "scrolls.db"),
			path.Join(loc, "scrolls.db"),
		)
		if err != nil {
			fmt.Printf("failed to move scrolls DB to new location. Err: %s\n", err.Error())
			return
		}

		settings.PersistChanges()
		fmt.Printf("Library location set to: %s\n", loc)
	},
}
