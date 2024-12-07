package cmd

import (
	"fmt"

	"github.com/mreliasen/scrolls-cli/internal/utils"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Returns the current installed version of Scrolls",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Scrolls CLI %s", utils.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
