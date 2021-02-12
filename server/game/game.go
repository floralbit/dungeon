package game

import (
	"time"

	"github.com/floralbit/dungeon/game/action"
	"github.com/floralbit/dungeon/game/data"
	"github.com/floralbit/dungeon/game/dungeon"
	"github.com/floralbit/dungeon/game/entity"
	"github.com/floralbit/dungeon/game/event"
	"github.com/floralbit/dungeon/game/event/network"
	"github.com/floralbit/dungeon/game/zone"
	"github.com/floralbit/dungeon/model"
	"github.com/google/uuid"
)

const tickLength = 100 // in ms
const eventBufferSize = 256

var startingZoneUUID = uuid.MustParse("10f8b073-cbd7-46b7-a6e3-9cbdf68a933f")

// In ...
var In = make(chan model.ClientEvent, eventBufferSize)

var zones = map[uuid.UUID]*zone.Zone{}

// Run ...
func Run() {
	event.Observers = append(event.Observers, network.NewObserver())

	zones = data.LoadZones()

	// build the dungeon!
	dungeonFloors := dungeon.BuildDungeon(zones[startingZoneUUID])
	for _, floor := range dungeonFloors {
		zones[floor.UUID] = floor
	}

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
		z.Update(dt)
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
	_, ok := entity.ActivePlayers[e.Sender.Account.UUID]
	if ok {
		return // player already logged in, TODO: handle gracefully ?
	}

	p := entity.NewPlayer(e.Sender) // TODO: pull from storage
	event.NotifyObservers(event.JoinEvent{Entity: p})
	p.Spawn(zones[startingZoneUUID])
}

func handleLeaveEvent(e model.ClientEvent) {
	p, ok := entity.ActivePlayers[e.Sender.Account.UUID]
	if !ok {
		return
	}
	p.Despawn()
}

func handleChatEvent(e model.ClientEvent) {
	p, ok := entity.ActivePlayers[e.Sender.Account.UUID]
	if !ok {
		return // ignore inactive players, TODO: send an error?
	}
	event.NotifyObservers(event.ChatEvent{Entity: p, Message: e.Chat.Message})
}

func handleMoveEvent(e model.ClientEvent) {
	p, ok := entity.ActivePlayers[e.Sender.Account.UUID]
	if !ok {
		return // ignore inactive players
	}
	p.QueuedAction = &action.MoveAction{
		Mover: p,
		X:     e.Move.X,
		Y:     e.Move.Y,
	}
}

func handleAttackEvent(e model.ClientEvent) {
	p, ok := entity.ActivePlayers[e.Sender.Account.UUID]
	if !ok {
		return // ignore inactive players
	}
	p.QueuedAction = &action.LightAttackAction{
		Attacker: p,
		X:        e.Attack.X,
		Y:        e.Attack.Y,
	}
}
