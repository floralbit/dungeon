package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/floralbit/dungeon/game"
	"github.com/floralbit/dungeon/model"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var sessionStore = sessions.NewCookieStore([]byte("todo-replace-with-real-env-key"))

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/game", handleGame)
	http.HandleFunc("/ws", handleWs)

	go game.Run()

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	err := authenticate(w, r, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fmt.Fprintf(w, "OK")
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	err := authenticate(w, r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fmt.Fprintf(w, "OK")
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, "dungeon")
	session.Values["UUID"] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleGame(w http.ResponseWriter, r *http.Request) {
	accountUUID, err := authenticated(w, r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	fmt.Fprintf(w, "Logged in with UUID %s", accountUUID.String())
}

// spawned in a goroutine by http
func handleWs(w http.ResponseWriter, r *http.Request) {
	// upgrade to websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// auth
	accountUUID, err := authenticated(w, r)
	if err != nil {
		ws.Close()
		return
	}

	fmt.Printf("%s connected\n", accountUUID.String())

	c := model.NewClient(ws, game.In)
	go c.HandleOutbound()
	c.HandleInbound()
}

func authenticate(w http.ResponseWriter, r *http.Request, newAccount bool) error {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		return errors.New("username or password not set")
	}

	var account *model.Account
	var err error
	if newAccount {
		account, err = model.Register(username, password)
		if err != nil {
			return err
		}
	} else {
		account, err = model.Login(username, password)
		if err != nil {
			return err
		}
	}

	session, _ := sessionStore.Get(r, "dungeon")

	session.Values["authenticated"] = true
	session.Values["UUID"] = account.UUID.String()
	session.Save(r, w)

	return nil
}

// returns account UUID if authenticated, errors if not
func authenticated(w http.ResponseWriter, r *http.Request) (*uuid.UUID, error) {
	session, _ := sessionStore.Get(r, "dungeon")

	rawUUID, ok := session.Values["UUID"]
	if !ok || rawUUID == "" {
		return nil, errors.New("not authenticated")
	}

	accountUUID, err := uuid.Parse(rawUUID.(string))
	if err != nil {
		return nil, err
	}

	return &accountUUID, nil
}
