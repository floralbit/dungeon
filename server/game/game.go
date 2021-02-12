package game

import (
	"github.com/floralbit/dungeon/game/event"
	"github.com/floralbit/dungeon/game/event/network"
	"time"

	"github.com/floralbit/dungeon/model"
)

const tickLength = 100 // in ms
const eventBufferSize = 256

// In ...
var In = make(chan model.ClientEvent, eventBufferSize)

var observers = []event.Observer{network.NewObserver()}

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

func update(dt float64) {
	for _, z := range zones {
		z.update(dt)
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
	case e.Attack != nil:
		handleAttackEvent(e)
	}
}

func handleJoinEvent(e model.ClientEvent) {
	_, ok := activePlayers[e.Sender.Account.UUID]
	if ok {
		return // player already logged in, TODO: handle gracefully ?
	}

	p := newPlayer(e.Sender) // TODO: pull from storage
	notifyObservers(event.JoinEvent{Entity: p})
	p.Spawn(startingZoneUUID)
}

func handleLeaveEvent(e model.ClientEvent) {
	p, ok := activePlayers[e.Sender.Account.UUID]
	if !ok {
		return
	}
	p.Despawn()
}

func handleChatEvent(e model.ClientEvent) {
	p, ok := activePlayers[e.Sender.Account.UUID]
	if !ok {
		return // ignore inactive players, TODO: send an error?
	}
	notifyObservers(event.ChatEvent{Entity: p, Message: e.Chat.Message})
}

func handleMoveEvent(e model.ClientEvent) {
	p, ok := activePlayers[e.Sender.Account.UUID]
	if !ok {
		return // ignore inactive players
	}
	p.queuedAction = &moveAction{
		Mover: p,
		X:     e.Move.X,
		Y:     e.Move.Y,
	}
}

func handleAttackEvent(e model.ClientEvent) {
	p, ok := activePlayers[e.Sender.Account.UUID]
	if !ok {
		return // ignore inactive players
	}

	p.queuedAction = &lightAttackAction{
		Attacker: p,
		X:        e.Attack.X,
		Y:        e.Attack.Y,
	}
}

func notifyObservers(e event.Event) {
	for _, o := range observers {
		o.Notify(e)
	}
}
