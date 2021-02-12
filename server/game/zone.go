package game

import (
	"github.com/floralbit/dungeon/game/model"
	"github.com/google/uuid"
)

type zone struct {
	UUID   uuid.UUID    `json:"uuid"`
	Name   string       `json:"name"`
	Width  int          `json:"width"`
	Height int          `json:"height"`
	Tiles  []model.Tile `json:"tiles"`

	Entities     map[uuid.UUID]model.Entity       `json:"entities"`
	WorldObjects map[uuid.UUID]*model.WorldObject `json:"world_objects"`
}

func (z *zone) GetUUID() uuid.UUID {
	return z.UUID
}

func (z *zone) GetDimensions() (int, int) {
	return z.Width, z.Height
}

func (z *zone) GetTile(x, y int) *model.Tile {
	if x < 0 || y < 0 || x >= z.Width || y >= z.Height {
		return nil
	}

	index := (z.Width * y) + x
	return &z.Tiles[index]
}

func (z *zone) GetEntities() (entities []model.Entity) {
	for _, e := range z.Entities {
		entities = append(entities, e)
	}
	return
}

func (z *zone) GetWorldObjects(x, y int) []*model.WorldObject {
	objs := []*model.WorldObject{}
	for _, obj := range z.WorldObjects {
		if obj.X == x && obj.Y == y {
			objs = append(objs, obj)
		}
	}
	return objs
}

func (z *zone) GetAllWorldObjects() []*model.WorldObject {
	objs := []*model.WorldObject{}
	for _, obj := range z.WorldObjects {
		objs = append(objs, obj)
	}
	return objs
}

func (z *zone) AddEntity(e model.Entity) {
	z.Entities[e.GetUUID()] = e
	e.SetZone(z)
}

func (z *zone) RemoveEntity(e model.Entity) {
	delete(z.Entities, e.GetUUID())
}

func (z *zone) update(dt float64) {
	actions := []model.Action{}
	for _, e := range z.Entities {
		if e.Tick() {
			a := e.Act()
			if a != nil {
				actions = append(actions, a)
			}
		}
	}

	for _, a := range actions {
		a.Execute()
	}
}
