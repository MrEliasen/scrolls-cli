package cmd

import (
	"fmt"
	"os"

	"github.com/mreliasen/scrolls-cli/internal/tui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(completionCmd)
	completionCmd.AddCommand(setBashCompletionCmd)
	completionCmd.AddCommand(setFishCompletionCmd)
	completionCmd.AddCommand(setZshCompletionCmd)
	completionCmd.AddCommand(setPowerShellCompletionCmd)
}

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate Scrolls CLI completion",
	Long: `Generate the autocompletion script for scrolls for the specified shell.
See each sub-command's help for details on how to use the generated script.`,
}

var setBashCompletionCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generate the autocompletion script for the bash shell.",
	Long: fmt.Sprintf(`
Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:
$ source <(scrolls completion bash)

To load completions for every new session, execute once:
%s:
$ scrolls completion bash > /etc/bash_completion.d/scrolls

%s:
$ scrolls completion bash > $(brew --prefix)/etc/bash_completion.d/scrolls

You will need to start a new shell for this setup to take effect.`,
		tui.HighlightStyle.Render("Linux"),
		tui.HighlightStyle.Render("macOS")),
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout)
	},
}

var setFishCompletionCmd = &cobra.Command{
	Use:   "fish",
	Short: "Generate the autocompletion script for the fish shell.",
	Long: `
Generate the autocompletion script for the fish shell.

To load completions in your current shell session:
$ scrolls completion fish | source

To load completions for every new session, execute once:
$ scrolls completion fish > ~/.config/fish/completions/scrolls.fish

You will need to start a new shell for this setup to take effect.`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenFishCompletion(os.Stdout, true)
	},
}

var setZshCompletionCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generate the autocompletion script for the zsh shell.",
	Long: fmt.Sprintf(`
Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:
$ source <(scrolls completion zsh)


To load completions for every new session, execute once:

%s:
$ scrolls completion zsh > "${fpath[1]}/_scrolls"

%s:
$ scrolls completion zsh > $(brew --prefix)/share/zsh/site-functions/_scrolls

You will need to start a new shell for this setup to take effect.`,
		tui.HighlightStyle.Render("Linux"),
		tui.HighlightStyle.Render("macOS")),
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenZshCompletion(os.Stdout)
	},
}

var setPowerShellCompletionCmd = &cobra.Command{
	Use:   "powershell",
	Short: "Generate the autocompletion script for powershell.",
	Long: `
Generate the autocompletion script for powershell.

To load completions in your current shell session:

scrolls completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
	},
}
