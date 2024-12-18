package flags

import "github.com/spf13/cobra"

var forceMigrate bool

func AddForceMigrateFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&forceMigrate, "force-migrate", false, "")
	cmd.PersistentFlags().MarkHidden("force-migrate")
}

func ForceMigrate() bool {
	return forceMigrate
}
