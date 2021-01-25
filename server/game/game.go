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
	// TODO: load from storage if exists
	player := newPlayer(e.Sender.Account.UUID, e.Sender.Account.Username)
	res := model.ServerEvent{
		Join: &model.ServerJoinEvent{
			Data: *player,
			From: e.Sender.Account.Username,
		},
	}

	// for now just pass to all clients
	for _, client := range model.ConnToClient {
		if client.Account.UUID != e.Sender.Account.UUID {
			client.In <- res
		}
	}

	// tell joiner they joined
	selfJoinRes := model.ServerEvent{
		Join: &model.ServerJoinEvent{
			Data: *player,
			From: e.Sender.Account.Username,
			You:  true,
		},
	}
	e.Sender.In <- selfJoinRes

	// tell joiner the zone data
	zoneRes := model.ServerEvent{
		Zone: &model.ServerZoneEvent{
			Data: Zones[player.Zone],
		},
	}
	e.Sender.In <- zoneRes
}

func handleLeaveEvent(e model.ClientEvent) {
	// TODO: save to storage
	delete(ActivePlayers, e.Sender.Account.UUID)

	res := model.ServerEvent{
		Leave: &model.ServerLeaveEvent{
			From: e.Sender.Account.Username,
		},
	}

	// for now just pass to all clients
	for _, client := range model.ConnToClient {
		client.In <- res
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

func handleMoveEvent(e model.ClientEvent) {
	// TODO: populate with UUIDs instead of username for lookup
	res := model.ServerEvent{
		Move: &model.ServerMoveEvent{
			X:    e.Move.X,
			Y:    e.Move.Y,
			From: e.Sender.Account.Username,
		},
	}

	for _, client := range model.ConnToClient {
		if client.Account.UUID != e.Sender.Account.UUID {
			client.In <- res
		}
	}

	senderRes := model.ServerEvent{
		Move: &model.ServerMoveEvent{
			X:    e.Move.X,
			Y:    e.Move.Y,
			From: e.Sender.Account.Username,
			You:  true,
		},
	}
	e.Sender.In <- senderRes
}
