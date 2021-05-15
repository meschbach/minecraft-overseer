package cmd

import "github.com/spf13/cobra"

func NewOverseerCLI() *cobra.Command {
	root := &cobra.Command {
		Use: "minecraft-overseer",
		Short: "Adapter to making Minecraft meet modern operational practices",
	}
	root.AddCommand(newServerCommands())
	return root
}