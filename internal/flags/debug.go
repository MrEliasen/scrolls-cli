package flags

import "github.com/spf13/cobra"

var debugFlag bool

func AddDebugFlag(cmd *cobra.Command) {
	desc := "If set, will dump all HTTP request and more."
	cmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, desc)
	cmd.PersistentFlags().MarkHidden("debug")
}

func Debug() bool {
	return debugFlag
}
