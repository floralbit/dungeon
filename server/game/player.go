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
	rollPlayerStats(&p)

	activePlayers[p.UUID] = &p
	return &p
}

// for new players
func rollPlayerStats(e *entity) {
	e.Stats.Level = 1

	// use 3d6 for stats
	r := roll{6, 1, 0} // 3d6 + 0
	e.Stats.Strength = r.roll()
	e.Stats.Dexterity = r.roll()
	e.Stats.Constitution = r.roll()
	e.Stats.Intelligence = r.roll()
	e.Stats.Wisdom = r.roll()
	e.Stats.Charisma = r.roll()

	// hit dice for players is a d8, so HP = 1d8 + CON
	e.Stats.HP = 0 // reset, this could be called on death
	for e.Stats.HP <= 0 {
		e.Stats.HP = roll{8, 1, modifier(e.Stats.Constitution)}.roll()
	}
}
