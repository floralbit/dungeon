package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/floralbit/dungeon/game"
	"github.com/floralbit/dungeon/model"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	// pages
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/game", handleGame)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../public/static"))))

	// endpoints
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/logout", handleLogout)

	// websocket
	http.HandleFunc("/ws", handleWs)

	go game.Run() // kick off gameloop

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// spawned in a goroutine by http
func handleWs(w http.ResponseWriter, r *http.Request) {
	// upgrade to websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// auth
	account, err := authenticated(w, r)
	if err != nil {
		ws.Close()
		return
	}

	fmt.Printf("%s connected\n", account.UUID.String())

	c := model.NewClient(ws, game.In, account)
	go c.HandleOutbound()
	c.HandleInbound()
}
