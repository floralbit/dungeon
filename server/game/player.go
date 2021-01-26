package game

import (
	"github.com/floralbit/dungeon/model"
	"github.com/google/uuid"
)

const warriorTileId = 21

var activePlayers = map[uuid.UUID]*entity{}

func newPlayer(client *model.Client) *entity {
	p := entity{
		UUID: client.Account.UUID,
		Name: client.Account.Username,
		Tile: warriorTileId,
		Type: entityTypePlayer,

		X: 24, Y: 18, // TODO: pull spawn from map data

		client: client,
	}
	activePlayers[p.UUID] = &p
	return &p
}

func (e *entity) leave() {
	e.zone.removeEntity(e)
	if e.Type == entityTypePlayer {
		delete(activePlayers, e.UUID)
	}
}