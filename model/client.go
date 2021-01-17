package model

import (
	"errors"
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

// Close ...
func (c *Client) Close() {
	// TODO: notify gameloop that client dropped
	c.Conn.Close()
	delete(ConnToClient, c.Conn)
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

		switch {
		case e.Register != nil:
			account, err := Register(e.Register.Username, e.Register.Password)
			if err != nil {
				c.SendError(err)
				continue
			}

			c.Account = account
			c.In <- ServerEvent{
				Register: &ServerRegisterEvent{
					Ok: true,
				},
			}

			c.Out <- e // send event to gameloop so new player can be made

		case e.Login != nil:
			account, err := Login(e.Login.Username, e.Login.Password)
			if err != nil {
				c.SendError(err)
				continue
			}

			c.Account = account
			c.In <- ServerEvent{
				Login: &ServerLoginEvent{
					Ok: true,
				},
			}

			c.Out <- e // send event to gameloop so new player can be spawned

		case e.Chat != nil:
			if c.Account == nil {
				c.SendError(errors.New("must be logged in to chat"))
				continue
			}

			c.Out <- e // send event to gameloop
		}

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
