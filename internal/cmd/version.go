package cmd

import (
	"fmt"
	"log"

	semver "github.com/hashicorp/go-version"
	"github.com/mreliasen/scrolls-cli/internal/scrolls"
	"github.com/mreliasen/scrolls-cli/internal/tui"
	"github.com/mreliasen/scrolls-cli/internal/utils"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Returns the current installed version of Scrolls",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Scrolls CLI %s\n", utils.Version)

		c, err := scrolls.New()
		if err != nil {
			log.Fatal("failed to initiate scroll-cli")
			return
		}

		msg, err := tui.NewSpinner("getting latest version info.", func() (string, error) {
			latest, err := c.Version.GetLatestReleaseVersion()
			if err != nil {
				return "", fmt.Errorf("error fetching latest version: %w", err)
			}

			parsedVersion, err := semver.NewVersion(utils.Version)
			if err != nil {
				return "", fmt.Errorf("error parsing current version: %w", err)
			}

			parsedLatest, err := semver.NewVersion(latest.Version)
			if err != nil {
				return "", fmt.Errorf("error parsing latest version: %w", err)
			}

			updateAvailable := parsedVersion.LessThan(parsedLatest)
			out := ""

			if updateAvailable {
				out += fmt.Sprintf("An update is available: %s\n", parsedLatest.String())
				out += fmt.Sprintf("Run %s to update to the latest version.\n", tui.HighlightStyle.Render("scrolls update"))
			} else {
				out += "You are running the latest version."
			}

			return out, nil
		})
		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Println(msg)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
