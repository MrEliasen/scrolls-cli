package cmd

import (
	"fmt"
	"os"

	"github.com/mreliasen/scrolls-cli/internal/flags"
	"github.com/mreliasen/scrolls-cli/internal/scrolls"
	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read <scroll-name>",
	Short: "Writes the scroll's content to stdout",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c, err := scrolls.New()
		if err != nil {
			fmt.Printf("failed to initiate scroll-cli")
			return
		}

		s, err := c.Files.GetScroll(name)
		if err != nil {
			if flags.Debug() {
				fmt.Printf("%+v\n", err)
			}

			fmt.Printf("failed to retrieve scroll.")
			return
		}

		if s.Type == "plain-text" {
			fmt.Printf("unable to cast a plain text scroll")
			return
		}

		os.Stdout.Write([]byte(s.Body()))
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}
