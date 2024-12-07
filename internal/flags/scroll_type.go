package flags

import "github.com/spf13/cobra"

var scrollTypeFlag string

func AddScrollTypeFlag(cmd *cobra.Command) {
	desc := "If set, will list all scrolls with the given file type."
	cmd.PersistentFlags().StringVarP(&scrollTypeFlag, "type", "t", "all", desc)
}

func ScrollType() string {
	return scrollTypeFlag
}
