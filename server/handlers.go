package main

import (
	"net/http"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../public/index.html")
}

func handleGame(w http.ResponseWriter, r *http.Request) {
	_, err := authenticated(w, r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.ServeFile(w, r, "../public/game.html")
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	err := authenticate(w, r, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	err := authenticate(w, r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, "dungeon")
	session.Values["UUID"] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
