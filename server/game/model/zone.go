package model

import "github.com/google/uuid"

type Tile struct {
	ID    int    `json:"id"`
	Solid bool   `json:"solid"`
	Name  string `json:"name"`
}

type Zone interface {
	GetUUID() uuid.UUID
	GetDimensions() (int, int)

	GetTile(x, y int) *Tile
	GetEntities() []Entity
	AddEntity(Entity)
	RemoveEntity(Entity)

	GetAllWorldObjects() []*WorldObject
	GetWorldObjects(x, y int) []*WorldObject
}
