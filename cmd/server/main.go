package main

import (
	"context"
	"fmt"
	"github.com/meschbach/go-junk-bucket/sub"
	"github.com/meschbach/minecraft-overseer/internal/junk"
	"github.com/meschbach/minecraft-overseer/internal/mc"
	"github.com/meschbach/minecraft-overseer/internal/mc/events"
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

	stdoutChannel := make(chan string, 16)
	echoStdoutChannel := make(chan string, 16)
	gameEventsInput := make(chan string, 16)
	stdout := &junk.StringBroadcast{
		Input: stdoutChannel,
		Out:   []chan<- string{echoStdoutChannel, gameEventsInput},
	}
	stderr := make(chan string, 16)
	stdin := make(chan string, 16)
	cmd := sub.NewSubcommand("java", []string{
		"-Dlog4j2.formatMsgNoLookups=true", "-Dlog4j.configurationFile=/log4j.xml",
		"-jar", "minecraft_server.jar",
		"--nogui"})
	go func() {
		fmt.Println("<<stderr initialized>>")
		for msg := range stderr {
			fmt.Fprintf(os.Stderr, "<<stderr>> %s\n", msg)
		}
	}()

	//standard output
	go stdout.RunLoop()

	go func() {
		fmt.Println("<<stdout initialized>>")
		for msg := range echoStdoutChannel {
			fmt.Printf("<<stdout>> %s\n", msg)
		}
	}()

	reactor := mc.StartReactor(gameEventsInput, stdin)
	reactor.PendingOperations <- &mc.WaitForStart{}
	reactor.PendingOperations <- &mc.EnsureUserOperators{Users: runtimeConfig.operators}
	reactor.PendingOperations <- &mc.EnsureWhitelistAdd{Users: runtimeConfig.users}

	gameChannel := make(chan events.LogEntry, 16)
	completedGamesChannel := reactor.Logs.Add(gameChannel)
	go func() {
		defer completedGamesChannel()
		for entry := range gameChannel {
			switch e := entry.(type) {
			case *events.UserJoinedEntry:
				fmt.Printf("User %q joined", e.User)
			case *events.UserLeftEvent:
				fmt.Printf("User %q left", e.User)
			default:
			}
		}
	}()

	err = cmd.Interact(stdin, stdoutChannel, stderr)
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
