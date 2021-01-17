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

	Out chan<- ClientEvent // to gameloop
	In  chan ServerEvent   // from gameLoop
}

// NewClient ...
func NewClient(conn *websocket.Conn, outChan chan<- ClientEvent) *Client {
	c := &Client{
		Conn: conn,
		Out:  outChan,
		In:   make(chan ServerEvent, eventBufferSize),
	}
	ConnToClient[conn] = c
	return c
}

// HandleInbound runs in websocket handler's goroutine (per conn)
func (c *Client) HandleInbound() {
	for {
		var e ClientEvent
		err := c.Conn.ReadJSON(&e)
		if err != nil {
			log.Printf("error: %v", err)
			delete(ConnToClient, c.Conn)
			break
		}

		c.Out <- e // send event to gameloop
	}
}

// HandleOutbound runs in its own goroutine too
func (c *Client) HandleOutbound() {
	for e := range c.In {
		err := c.Conn.WriteJSON(e)
		if err != nil {
			log.Printf("error: %v", err)
			c.Conn.Close()
			delete(ConnToClient, c.Conn)
		}
	}
}
