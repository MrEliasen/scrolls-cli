package flags

import "github.com/spf13/cobra"

var templateFlag bool

func AddTemplateFlag(cmd *cobra.Command) {
	desc := "If set, will populate the scroll with a template for the selected file type."
	cmd.PersistentFlags().BoolVarP(&templateFlag, "template", "t", false, desc)
}

func Template() bool {
	return templateFlag
}
