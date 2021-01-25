package game

import "github.com/google/uuid"

const warriorTile = 21

// Player ...
type Player struct {
	UUID uuid.UUID
	Name string
	Tile int

	X, Y int
	Zone string
}

// ActivePlayers ...
var ActivePlayers = map[uuid.UUID]*Player{}

func newPlayer(UUID uuid.UUID, name string) *Player {
	p := Player{
		UUID: UUID,
		Name: name,
		Tile: warriorTile,

		X:    24, // TODO: pull spawn from map data
		Y:    18,
		Zone: startingZone,
	}
	ActivePlayers[UUID] = &p
	return &p
}
