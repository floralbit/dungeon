package game

import (
	"time"

	"github.com/floralbit/dungeon/model"
)

const tickLength = 100 // in ms
const eventBufferSize = 256

// In ...
var In = make(chan model.ClientEvent, eventBufferSize)

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
			// TODO: do something about this nasty dispatch, will be a pain to maintain
			// TODO: notify others when someone joins (login) or leaves
			switch {
			case e.Join != nil:
				handleJoinEvent(e)
			case e.Leave != nil:
				handleLeaveEvent(e)
			case e.Chat != nil:
				handleChatEvent(e)
			case e.Move != nil:
				handleMoveEvent(e)
			}
		default:
			return // In is empty
		}
	}
}

func update(dt float64) {
	for _, z := range zones {
		z.update(dt)
	}
}

func handleJoinEvent(e model.ClientEvent) {
	_, ok := activePlayers[e.Sender.Account.UUID]
	if ok {
		return // player already logged in, TODO: handle gracefully ?
	}

	p := newPlayer(e.Sender)                   // TODO: pull from storage
	p.Send(newServerMessageEvent(motd, false)) // send message of the day
	p.Spawn(startingZoneUUID)
}

func handleLeaveEvent(e model.ClientEvent) {
	p, ok := activePlayers[e.Sender.Account.UUID]
	if !ok {
		return
	}
	p.Despawn(false)
}

func handleChatEvent(e model.ClientEvent) {
	p, ok := activePlayers[e.Sender.Account.UUID]
	if !ok {
		return // ignore inactive players, TODO: send an error?
	}
	p.Data().zone.send(newChatEvent(p.Data(), e.Chat.Message))
}

func handleMoveEvent(e model.ClientEvent) {
	p, ok := activePlayers[e.Sender.Account.UUID]
	if !ok {
		return // ignore inactive players
	}
	p.Move(e.Move.X, e.Move.Y)
}
