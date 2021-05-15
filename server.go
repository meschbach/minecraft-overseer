package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/meschbach/minecraft-overseer/discord"
	wui "github.com/meschbach/minecraft-overseer/wui"
	ws "github.com/meschbach/minecraft-overseer/ws"
)

var (
	addr    = flag.String("addr", "127.0.0.1:8080", "http service address")
	discordToken = flag.String("discord", "", "discord token to active discord")
)

func internalError(ws *websocket.Conn, msg string, err error) {
	log.Println(msg, err)
	ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}

func main() {
	initCtx := context.Background()
	flag.Parse()

	hub := ws.NewHub()
	if discordToken != nil {
		fmt.Printf("Creating discord client\n")
		_, err := discord.NewDiscordClient(initCtx, *discordToken)
		if err != nil {
			panic(err)
		}
	}

	go hub.Run()

	http.HandleFunc("/", wui.ServeWUI)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ProcessClient(hub, w, r)
	})
	//TODO: Binds to 8080 no matter what
	log.Fatal(http.ListenAndServe(*addr, nil))
}
