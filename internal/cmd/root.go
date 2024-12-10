package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/mreliasen/scrolls-cli/internal/scrolls"
	"github.com/mreliasen/scrolls-cli/internal/settings"
	"github.com/mreliasen/scrolls-cli/internal/tui"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scrolls",
	Short: "",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {
		config, err := settings.LoadSettings()
		if err != nil {
			return
		}

		if time.Now().Unix() <= config.GetLastUpdateCheck()+int64(12*60*60) {
			return
		}

		u, err := scrolls.New()
		if err != nil {
			return
		}

		curr, latest, update, err := u.Version.CheckForUpdates(false)
		if err != nil {
			return
		}

		if update {
			fmt.Printf("\n\nHeads up! There is a newer version of %s available.\n", tui.HighlightStyle.Render("Turso CLI"))
			fmt.Printf("You're currently on version %s, the latest available version is %s.\n", tui.HighlightStyle.Render(curr), tui.HighlightStyle.Render(latest))
			fmt.Printf("To update:\n\n\t%s\n\n", tui.HighlightStyle.Render("scrolls update"))
		}
	}
}
