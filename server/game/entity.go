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
	UUID uuid.UUID `json:"uuid"`
	Name string `json:"name"`
	Tile int `json:"tile"` // representing tile
	Type entityType `json:"type"`

	X int `json:"x"`
	Y int `json:"y"`

	zone *zone `json:"-"`
	client *model.Client `json:"-"`
}

type worldObject struct {
	UUID uuid.UUID `json:"uuiud"`
	Name string `json:"name"`
	Tile int `json:"tile"` // representing tile
	X int `json:"int"`
	Y int `json:"int"`
}
