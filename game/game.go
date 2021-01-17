package game

import (
	"time"

	"github.com/floralbit/dungeonserv/model"
)

const tickLength = 100 // in ms
const eventBufferSize = 256

// In ...
var In = make(chan model.Event, eventBufferSize)

// Out ...
var Out = make(chan model.Event, eventBufferSize)

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
		case msg := <-In:
			handleMessage(msg)
		default:
			return // In is empty
		}
	}
}

func update(dt float64) {
	// do game logic updates (monsters, whatever)
}

func handleMessage(e model.Event) {
	// TODO: do something about this nasty dispatch, will be a pain to maintain
	switch {
	case e.ChatMessage != nil:
		handleChatMessage(e.ChatMessage)
	}
}

func handleChatMessage(msg *model.ChatMessage) {
	res := model.Event{
		ChatMessage: &model.ChatMessage{
			Data: msg.Data,
		},
	}

	// for now just pass to all clients
	for _, client := range model.ConnToClient {
		client.In <- res
	}
}
