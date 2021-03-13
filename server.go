package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	addr    = flag.String("addr", "127.0.0.1:8080", "http service address")
)

func internalError(ws *websocket.Conn, msg string, err error) {
	log.Println(msg, err)
	ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("usage: minecraft-oversee <http-address>")
	}
	hub := newHub()
	go hub.run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	//TODO: Binds to 8080 no matter what
	log.Fatal(http.ListenAndServe(*addr, nil))
}
