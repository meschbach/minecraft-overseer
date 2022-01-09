package cmd

import "github.com/spf13/cobra"

func NewOverseerCLI() *cobra.Command {
	root := &cobra.Command{
		Use:           "overseer",
		Short:         "Adapter to making Minecraft meet modern operational practices",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.AddCommand(newServerCommands())
	root.AddCommand(newInitCommands())
	return root
}