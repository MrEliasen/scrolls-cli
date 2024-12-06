package cmd

import (
	"fmt"

	"github.com/mreliasen/scrolls-cli/internal/scrolls"
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

		c.Version.CheckForUpdates()
		return nil
	},
}
