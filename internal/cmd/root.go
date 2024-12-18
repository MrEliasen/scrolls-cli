package cmd

import (
	"fmt"
	"os"
	"time"

	semver "github.com/hashicorp/go-version"
	"github.com/mreliasen/scrolls-cli/internal/flags"
	"github.com/mreliasen/scrolls-cli/internal/library"
	"github.com/mreliasen/scrolls-cli/internal/scrolls"
	"github.com/mreliasen/scrolls-cli/internal/settings"
	"github.com/mreliasen/scrolls-cli/internal/tui"
	"github.com/mreliasen/scrolls-cli/internal/utils"
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
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		config, err := settings.LoadSettings()
		if err != nil {
			return
		}

		runMigrations := flags.ForceMigrate()
		ver := config.GetMigrationVersion()
		if ver == "" {
			ver = "v0.0.0"
		}

		if !runMigrations {

			parsedCLIVersion, err := semver.NewVersion(utils.Version)
			if err != nil {
				return
			}

			parsedConfigVersion, err := semver.NewVersion(ver)
			if err != nil {
				return
			}

			runMigrations = parsedConfigVersion.LessThan(parsedCLIVersion)
		}

		if runMigrations {
			lib, err := library.LoadLibrary()
			if err != nil {
				if flags.Debug() {
					fmt.Fprintf(os.Stderr, "failed to run migrations: %s\n", err.Error())
				}
				return
			}

			// first migration? also migrate scrolls to db
			if ver == "v0.0.0" {
				err := lib.MigrateScrolls(config.GetLibrary())
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to migrate scrolls to db\n")
					fmt.Fprintf(os.Stderr, "run with --debug to see more details.\n")

					if flags.Debug() {
						fmt.Fprintf(os.Stderr, "migration error: %s\n\n", err.Error())
					}

					return
				} else {
					fmt.Println("Scrolls have been migrated to SQLite, the old format scrolls still exists but are no longer in use.")
					fmt.Println("Kept purely \"just in case\" anything went wrong. You can find then in the old library path.")
					fmt.Printf("Old library path: %s\n\n", config.GetLibrary())
				}
			}

			config.SetMigrationVersion(utils.Version)
			config.PersistChanges()
		}
	}

	rootCmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {
		config, err := settings.LoadSettings()
		if err != nil {
			return
		}

		// close the open db
		lib, err := library.LoadLibrary()
		if err == nil {
			lib.Close()
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
			fmt.Printf("\nHeads up! There is a newer version of %s available.\n", tui.HighlightStyle.Render("Scrolls"))
			fmt.Printf("You're currently on version %s, the latest available version is %s.\n", tui.HighlightStyle.Render(curr), tui.HighlightStyle.Render(latest))
			fmt.Printf("To update:\t%s\n\n", tui.HighlightStyle.Render("scrolls update"))
		}
	}

	flags.AddDebugFlag(rootCmd)
	flags.AddResetConfigFlag(rootCmd)
	flags.AddForceMigrateFlag(rootCmd)
}
