package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/floralbit/dungeonserv/game"
	"github.com/floralbit/dungeonserv/model"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleWs)

	go game.Run()

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// spawned in a goroutine by http
func handleWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got ws connection...")

	// upgrade to websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	c := model.NewClient(ws, game.In)
	go c.HandleOutbound()
	c.HandleInbound()
}
