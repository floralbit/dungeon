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
	// do game logic updates (monsters, whatever)
}

func handleJoinEvent(e model.ClientEvent) {
	p := newPlayer(e.Sender) // TODO: pull from storage
	// TODO: dispatch to zone add entity method
	zones[startingZoneUUID].Entities[p.UUID] = p // add player to starting zone
	p.zone = zones[startingZoneUUID]
	zones[startingZoneUUID].send(newSpawnEvent(p)) // notify zone player has joined
	p.send(newZoneLoadEvent(p.zone)) // send player zone data
}

func handleLeaveEvent(e model.ClientEvent) {
	p := activePlayers[e.Sender.Account.UUID]
	// TODO: dispatch to zone remove entity method
	delete(p.zone.Entities, p.UUID)
	delete(activePlayers, p.UUID)
	p.zone.send(newDespawnEvent(p))
}

func handleChatEvent(e model.ClientEvent) {
	p := activePlayers[e.Sender.Account.UUID]
	p.zone.send(newChatEvent(p, e.Chat.Message))
}

func handleMoveEvent(e model.ClientEvent) {
	// TODO: pass to zone to handle for collision detection
	p := activePlayers[e.Sender.Account.UUID]
	p.zone.send(newMoveEvent(p, e.Move.X, e.Move.Y))
}
