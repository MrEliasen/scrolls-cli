package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mreliasen/scrolls-cli/internal/scrolls"
	"github.com/mreliasen/scrolls-cli/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the Scrolls CLI to the latest version",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := scrolls.New()
		if err != nil {
			return fmt.Errorf("failed to initiate scrolls cli: %w", err)
		}

		latest, err := c.Version.GetLatestRelease()
		if err != nil {
			return fmt.Errorf("failed to get version information: %w", err)
		}

		fmt.Printf("Current version: %s, latest version: %s\n", utils.Version, latest.Version)
		if utils.Version >= latest.Version {
			fmt.Printf("version %s is already latest\n", utils.Version)
			return nil
		}

		return Update()
	},
}

func Update() error {
	command := exec.Command("sh", "-c", "curl -sSfL \"https://get.scrolls.sh/install.sh\" | sh")
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		return fmt.Errorf("failed to execute update command: %w", err)
	}

	return nil
}
