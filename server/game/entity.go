package game

import (
	"github.com/floralbit/dungeon/model"
	"github.com/google/uuid"
)

type entityType string

const (
	entityTypePlayer = "player"
)

type entity struct {
	UUID uuid.UUID  `json:"uuid"`
	Name string     `json:"name"`
	Tile int        `json:"tile"` // representing tile
	Type entityType `json:"type"`

	X int `json:"x"`
	Y int `json:"y"`

	zone   *zone         `json:"-"`
	client *model.Client `json:"-"`
}

func (e *entity) move(x, y int) {
	t := e.zone.getTile(x, y)
	if t == nil {
		return // edge of map, don't move
	}

	if t.Solid {
		return
	}

	e.zone.send(newMoveEvent(e, x, y))
}
