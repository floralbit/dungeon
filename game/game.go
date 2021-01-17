package game

import (
	"time"

	"github.com/floralbit/dungeonserv/model"
)

const tickLength = 100 // in ms
const eventBufferSize = 256

// In ...
var In = make(chan model.ClientEvent, eventBufferSize)

// Out ...
var Out = make(chan model.ServerEvent, eventBufferSize)

// Run ...
func Run() {
	ticker := time.NewTicker(tickLength * time.Millisecond)
	lastTime := time.Now()

	for now := range ticker.C {
		dt := now.Sub(lastTime).Seconds()
		lastTime = now

		processEvents(dt)
		update(dt)
	}
}

func processEvents(dt float64) {
	// do we need a timeout timer in here too? we might starve the gameloop
	// if the inbounds are so fast that we never drain the channel
	for {
		select {
		case e := <-In:
			handleEvent(e)
		default:
			return // In is empty
		}
	}
}

func update(dt float64) {
	// do game logic updates (monsters, whatever)
}

func handleEvent(e model.ClientEvent) {
	// TODO: do something about this nasty dispatch, will be a pain to maintain
	switch {
	case e.Chat != nil:
		handleChatEvent(e)
	}
}

func handleChatEvent(e model.ClientEvent) {
	res := model.ServerEvent{
		Chat: &model.ServerChatEvent{
			Message: e.Chat.Message,
			From:    e.Sender.Account.Username,
		},
	}

	// for now just pass to all clients
	for _, client := range model.ConnToClient {
		client.In <- res
	}
}
