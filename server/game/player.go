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

//// Player ...
//type Player struct {
//	UUID uuid.UUID
//	Name string
//	Tile int
//
//	X, Y int
//	Zone string
//}
//
//// ActivePlayers ...
//var ActivePlayers = map[uuid.UUID]*Player{}
//
//func newPlayer(UUID uuid.UUID, name string) *Player {
//	p := Player{
//		UUID: UUID,
//		Name: name,
//		Tile: warriorTile,
//
//		X:    24, // TODO: pull spawn from map data
//		Y:    18,
//	}
//	ActivePlayers[UUID] = &p
//	return &p
//}
