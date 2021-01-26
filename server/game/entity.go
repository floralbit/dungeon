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

type tile struct {
	ID int `json:"id"`
	Solid bool `json:"solid"`
}

type zone struct {
	UUID uuid.UUID `json:"uuid"`
	Name string `json:"name"`
	Width int `json:"width"`
	Height int `json:"height"`
	Tiles []tile `json:"tiles"`

	Entities map[uuid.UUID]*entity `json:"entities"`
	WorldObjects map[uuid.UUID]*worldObject `json:"world_objects"`
}

type worldObject struct {
	UUID uuid.UUID `json:"uuiud"`
	Name string `json:"name"`
	Tile int `json:"tile"` // representing tile
	X int `json:"int"`
	Y int `json:"int"`
}
