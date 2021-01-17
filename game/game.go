package game

import (
	"fmt"
	"log"
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

func handleMessage(e model.Event) {
	// for now, just send to everyone connected lol
	res := model.Event{
		Data: e.Data,
	}

	fmt.Println(res)

	for _, client := range model.ConnToClient {
		err := client.Conn.WriteJSON(res)
		if err != nil {
			log.Printf("error: %v", err)
			client.Conn.Close()
			delete(model.ConnToClient, client.Conn)
		}
	}
}

func update(dt float64) {
	// do game logic updates (monsters, whatever)
}
