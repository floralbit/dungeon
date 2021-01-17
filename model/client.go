package model

import (
	"log"

	"github.com/gorilla/websocket"
)

const eventBufferSize = 256

// ConnToClient ...
var ConnToClient = map[*websocket.Conn]*Client{}

// Client ...
type Client struct {
	Conn    *websocket.Conn
	Account *Account

	Out chan Event // from connection to gameloop
	In  chan Event // from game to connection
}

// NewClient ...
func NewClient(conn *websocket.Conn, gameChan chan Event) *Client {
	c := &Client{
		Conn: conn,
		Out:  gameChan,
		In:   make(chan Event, eventBufferSize),
	}
	ConnToClient[conn] = c
	return c
}

// HandleInputs runs in websocket handler's goroutine (per conn)
func (c *Client) HandleInputs() {
	for {
		var e Event
		err := c.Conn.ReadJSON(&e)
		if err != nil {
			log.Printf("error: %v", err)
			delete(ConnToClient, c.Conn)
			break
		}

		c.Out <- e // send event to gameloop
	}
}

// HandleOutputs runs in its own goroutine too
func (c *Client) HandleOutputs() {
	for e := range c.In {
		err := c.Conn.WriteJSON(e)
		if err != nil {
			log.Printf("error: %v", err)
			c.Conn.Close()
			delete(ConnToClient, c.Conn)
		}
	}
}
