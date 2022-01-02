package cmd

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"

	"github.com/meschbach/minecraft-overseer/discord"
	ws "github.com/meschbach/minecraft-overseer/ws"
	wui "github.com/meschbach/minecraft-overseer/wui"
	"github.com/spf13/cobra"
)

func internalError(ws *websocket.Conn, msg string, err error) {
	log.Println(msg, err)
	ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}

func RunServer(initCtx context.Context, opts *serverOpts) error {
	hub := ws.NewHub()
	if len(opts.discordToken) > 0 {
		fmt.Printf("Creating discord client\n")
		_, err := discord.NewDiscordClient(initCtx, opts.discordToken)
		if err != nil {
			return err
		}
	}

	go hub.Run()
	//Automated initialization and start
	go func() {
		start := &ws.StartMessage{}
		hub.InternalSend(start)
	}()

	http.HandleFunc("/", wui.ServeWUI)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ProcessClient(hub, w, r)
	})
	fmt.Printf("Starting webserver at %q\n", opts.httpAddress)
	return http.ListenAndServe(opts.httpAddress, nil)
}

type serverOpts struct {
	httpAddress  string
	discordToken string
}

func newServerCommands() *cobra.Command {
	opts := &serverOpts{}
	run := &cobra.Command{
		Use:   "server",
		Short: "Begins the Overseer service",
		//PreRunE: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			startup := cmd.Context()
			return RunServer(startup, opts)
		},
	}
	run.PersistentFlags().StringVar(&opts.httpAddress, "http-bind", "127.0.0.1:8080", "Port to bind webhost too")
	run.PersistentFlags().StringVar(&opts.discordToken, "discord-token", "", "Enables connecting to Discord")
	return run
}
