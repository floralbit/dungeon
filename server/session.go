package main

import (
	"errors"
	"net/http"

	"github.com/floralbit/dungeon/store"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

var sessionStore = sessions.NewCookieStore([]byte("todo-replace-with-real-env-key"))

func authenticate(w http.ResponseWriter, r *http.Request, newAccount bool) error {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		return errors.New("username or password not set")
	}

	var account *store.Account
	var err error
	if newAccount {
		account, err = store.RegisterAccount(username, password)
		if err != nil {
			return err
		}
	} else {
		account, err = store.LoginAccount(username, password)
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
func authenticated(w http.ResponseWriter, r *http.Request) (*store.Account, error) {
	session, _ := sessionStore.Get(r, "dungeon")

	rawUUID, ok := session.Values["UUID"]
	if !ok || rawUUID == "" {
		return nil, errors.New("not authenticated")
	}

	accountUUID, err := uuid.Parse(rawUUID.(string))
	if err != nil {
		return nil, err
	}

	account, err := store.GetAccount(accountUUID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.New("account not found")
	}

	return account, nil
}
