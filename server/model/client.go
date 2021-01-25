package model

import (
	"log"

	"github.com/floralbit/dungeon/store"
	"github.com/gorilla/websocket"
)

const eventBufferSize = 256

// ConnToClient ...
var ConnToClient = map[*websocket.Conn]*Client{}

// Client ...
type Client struct {
	Conn    *websocket.Conn
	Account *store.Account

	Out chan<- ClientEvent // to gameloop
	In  chan ServerEvent   // from gameLoop
}

// NewClient ...
func NewClient(conn *websocket.Conn, outChan chan<- ClientEvent, account *store.Account) *Client {
	c := &Client{
		Conn:    conn,
		Out:     outChan,
		In:      make(chan ServerEvent, eventBufferSize),
		Account: account,
	}
	ConnToClient[conn] = c

	// notify join
	outChan <- ClientEvent{
		Sender: c,
		Join: &ClientJoinEvent{
			Ok: true,
		},
	}

	return c
}

// Close ...
func (c *Client) Close() {
	c.Conn.Close()
	delete(ConnToClient, c.Conn)

	// let game know they've bailed
	c.Out <- ClientEvent{
		Sender: c,
		Leave: &ClientLeaveEvent{
			Ok: true,
		},
	}
}

// HandleInbound runs in websocket handler's goroutine (per conn)
func (c *Client) HandleInbound() {
	for {
		var e ClientEvent
		err := c.Conn.ReadJSON(&e)
		if err != nil {
			log.Printf("error: %v", err)
			c.Close()
			return
		}

		e.Sender = c // label sender
		c.Out <- e   // send event to game loop
	}
}

// HandleOutbound runs in its own goroutine too
func (c *Client) HandleOutbound() {
	for e := range c.In {
		err := c.Conn.WriteJSON(e)
		if err != nil {
			log.Printf("error: %v", err)
			c.Close()
			return
		}
	}
}

// SendError ...
func (c *Client) SendError(err error) {
	c.In <- ServerEvent{
		Error: &ServerErrorEvent{
			Message: err.Error(),
		},
	}
}
