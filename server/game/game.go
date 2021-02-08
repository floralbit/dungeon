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

	for {
		select {
		case now := <-ticker.C:
			dt := now.Sub(lastTime).Seconds()
			lastTime = now
			update(dt)
		case e := <-In:
			processEvent(e)
		}
	}
}

func processEvent(e model.ClientEvent) {
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
