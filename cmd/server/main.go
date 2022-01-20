package main

import (
	"context"
	"fmt"
	"github.com/meschbach/minecraft-overseer/internal/discord"
	"github.com/meschbach/minecraft-overseer/internal/mc"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	root := &cobra.Command{
		Use:           "overseer-server",
		Short:         "Overseer instance specialized for server environments",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.AddCommand(newRunCommand())
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func RunProgram(initCtx context.Context, opts *serverOpts) error {
	runtimeConfig, err := initV2(initCtx, opts.fs.configFile, opts.fs.gameDir)
	if err != nil {
		return err
	}
	fmt.Printf("Passed configuration: %v\n", runtimeConfig)

	game, err := mc.NewInstance(runtimeConfig.gameDirectory)
	if err != nil {
		return err
	}
	instance, err := game.PrepareRunning()
	if err != nil {
		return err
	}

	reactor := instance.Reactor
	reactor.PendingOperations <- &mc.WaitForStart{}
	reactor.PendingOperations <- &mc.EnsureUserOperators{Users: runtimeConfig.operators}
	reactor.PendingOperations <- &mc.EnsureWhitelistAdd{Users: runtimeConfig.users}

	for _, discordSpec := range runtimeConfig.discord {
		logger, err := discord.NewLogger(discordSpec.token, "dev-minecraft-overseer")
		if err != nil {
			return err
		}
		go logger.Ingest(reactor.Logs)
	}

	err = instance.Run()
	if err != nil {
		return err
	}
	return nil
}

type serverOpts struct {
	httpAddress  string
	discordToken string
	fs           struct {
		configFile string
		gameDir    string
	}
}

func newRunCommand() *cobra.Command {
	opts := &serverOpts{}
	run := &cobra.Command{
		Use:   "run",
		Short: "Runs the service",
		//PreRunE: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			startup := cmd.Context()
			return RunProgram(startup, opts)
		},
	}
	run.PersistentFlags().StringVar(&opts.httpAddress, "http-bind", "127.0.0.1:8080", "Port to bind webhost too")
	run.PersistentFlags().StringVar(&opts.discordToken, "discord-token", "", "Enables connecting to Discord")
	run.PersistentFlags().StringVarP(&opts.fs.gameDir, "game-dir", "d", "/mc/instance", "Game directory")
	run.PersistentFlags().StringVarP(&opts.fs.configFile, "config-file", "c", "/mc/config/manifest.json", "Configuration manifest for game")
	return run
}
