package store

import (
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // justify it lol
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

// Account ...
type Account struct {
	UUID           uuid.UUID
	Username       string
	HashedPassword []byte
}

// Init ...
func Init() {
	var err error
	db, err = sql.Open("sqlite3", "../database.db")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

// LoginAccount ...
func LoginAccount(username, password string) (*Account, error) {
	account, err := GetAccountByUsername(username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(account.HashedPassword, []byte(password)); err != nil {
		return nil, errors.New("Passwords do not match")
	}

	return account, nil
}

// RegisterAccount ...
func RegisterAccount(username, password string) (*Account, error) {
	account, _ := GetAccountByUsername(username)
	if account != nil {
		return nil, errors.New("account already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	accountUUID := uuid.New()

	statement, err := db.Prepare("INSERT INTO accounts (uuid, username, hashed_password) VALUES (?, ?, ?)")
	if err != nil {
		return nil, err
	}

	_, err = statement.Exec(accountUUID.String(), username, string(hashedPassword))
	if err != nil {
		return nil, err
	}
	statement.Close()

	return GetAccount(accountUUID)
}

// GetAccount ...
func GetAccount(accountUUID uuid.UUID) (*Account, error) {
	var rawUUID, username, hashedPassword string

	statement, err := db.Prepare("SELECT uuid, username, hashed_password FROM accounts WHERE uuid = ?")
	if err != nil {
		return nil, err
	}

	err = statement.QueryRow(accountUUID.String()).Scan(&rawUUID, &username, &hashedPassword)
	if err != nil {
		return nil, err
	}
	statement.Close()

	parsedUUID, err := uuid.Parse(rawUUID)
	if err != nil {
		return nil, err
	}

	a := &Account{
		UUID:           parsedUUID,
		Username:       username,
		HashedPassword: []byte(hashedPassword),
	}

	return a, nil
}

// GetAccountByUsername ...
func GetAccountByUsername(username string) (*Account, error) {
	var rawUUID, resUsername, hashedPassword string

	statement, err := db.Prepare("SELECT uuid, username, hashed_password FROM accounts WHERE username = ?")
	if err != nil {
		return nil, err
	}

	err = statement.QueryRow(username).Scan(&rawUUID, &resUsername, &hashedPassword)
	if err != nil {
		return nil, err
	}
	statement.Close()

	parsedUUID, err := uuid.Parse(rawUUID)
	if err != nil {
		return nil, err
	}

	a := &Account{
		UUID:           parsedUUID,
		Username:       resUsername,
		HashedPassword: []byte(hashedPassword),
	}

	return a, nil
}
