package model

import "github.com/gorilla/websocket"

// ConnToClient ...
var ConnToClient = map[*websocket.Conn]*Client{}

// Client ...
type Client struct {
	Conn *websocket.Conn
}
