package cmd

import (
	"fmt"
	"os"
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
	configGetCmd.AddCommand(GetLibraryCmd)
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

var GetLibraryCmd = &cobra.Command{
	Use:   "library",
	Short: "Get the path to where your scrolls are stored.",
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

		fmt.Printf("Library location: %s\n", e)
	},
}

var setLibraryCmd = &cobra.Command{
	Use:   "library <path>",
	Short: "Set the path to where you want to store you scrolls. Scrolls will be automatically moved to the new location",
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
					fmt.Println("failed to create the new library folder, as it did not already exist.")
					return
				}
			} else {
				if !stat.IsDir() {
					fmt.Println("failed to set new library location. The specified location is not a directory.")
					return
				}
			}
		}

		settings.PersistChanges()

		files, err := os.ReadDir(src)
		if err != nil {
			fmt.Println("failed to set new library location.")
			fmt.Println(err.Error())
			return
		}

		fmt.Println("migrating scrolls, please wait..")

		i := 0
		f := 0
		for _, entry := range files {
			if entry.IsDir() {
				continue
			}

			err = os.Rename(
				fmt.Sprintf("%s/%s", src, entry.Name()),
				fmt.Sprintf("%s/%s", loc, entry.Name()),
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to move scroll: %s\n", entry.Name())
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				f++
			} else {
				i++
			}
		}

		fmt.Printf("Library location set to: %s\n", loc)
		fmt.Printf("Scrolls Moved: %d successfully, %d failed\n", i, f)
	},
}
