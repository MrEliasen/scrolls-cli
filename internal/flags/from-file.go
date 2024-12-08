package flags

import "github.com/spf13/cobra"

var fromFileFlag string

func AddFromFileFlag(cmd *cobra.Command) {
	desc := "Copy the content of the specified file as template for the new scroll. This flag takes priority over --template."
	cmd.PersistentFlags().StringVarP(&fromFileFlag, "file", "f", "", desc)
}

func FromFile() string {
	return fromFileFlag
}
