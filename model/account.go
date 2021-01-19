package model

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Accounts ...
var Accounts = map[uuid.UUID]*Account{} // TODO: use a real db lol

// Account ...
type Account struct {
	UUID           uuid.UUID
	Username       string
	HashedPassword []byte
}

// Login ...
func Login(username, password string) (*Account, error) {
	account := getAccountByUsername(username)
	if account == nil {
		return nil, errors.New("No account found")
	}

	if err := bcrypt.CompareHashAndPassword(account.HashedPassword, []byte(password)); err != nil {
		return nil, errors.New("Passwords do not match")
	}

	return account, nil
}

// Register ...
func Register(username, password string) (*Account, error) {
	if acc := getAccountByUsername(username); acc != nil {
		return nil, errors.New("account already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	a := &Account{
		UUID:           uuid.New(),
		Username:       username,
		HashedPassword: hashedPassword,
	}
	Accounts[a.UUID] = a

	return a, nil
}

// GetAccountByUUID ...
func GetAccountByUUID(UUID uuid.UUID) *Account {
	return Accounts[UUID]
}

// TODO: replace with actual db query instead of linear lookup
func getAccountByUsername(username string) *Account {
	for _, account := range Accounts {
		if account.Username == username {
			return account
		}
	}
	return nil
}
