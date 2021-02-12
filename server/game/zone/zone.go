package zone

import (
	"github.com/floralbit/dungeon/game/model"
	"github.com/google/uuid"
)

type Zone struct {
	UUID   uuid.UUID    `json:"uuid"`
	Name   string       `json:"name"`
	Width  int          `json:"width"`
	Height int          `json:"height"`
	Tiles  []model.Tile `json:"tiles"`

	Entities     map[uuid.UUID]model.Entity       `json:"entities"`
	WorldObjects map[uuid.UUID]*model.WorldObject `json:"world_objects"`
}

func (z *Zone) GetUUID() uuid.UUID {
	return z.UUID
}

func (z *Zone) GetDimensions() (int, int) {
	return z.Width, z.Height
}

func (z *Zone) GetTile(x, y int) *model.Tile {
	if x < 0 || y < 0 || x >= z.Width || y >= z.Height {
		return nil
	}

	index := (z.Width * y) + x
	return &z.Tiles[index]
}

func (z *Zone) GetEntities() (entities []model.Entity) {
	for _, e := range z.Entities {
		entities = append(entities, e)
	}
	return
}

func (z *Zone) GetWorldObjects(x, y int) []*model.WorldObject {
	objs := []*model.WorldObject{}
	for _, obj := range z.WorldObjects {
		if obj.X == x && obj.Y == y {
			objs = append(objs, obj)
		}
	}
	return objs
}

func (z *Zone) GetAllWorldObjects() []*model.WorldObject {
	objs := []*model.WorldObject{}
	for _, obj := range z.WorldObjects {
		objs = append(objs, obj)
	}
	return objs
}

func (z *Zone) AddEntity(e model.Entity) {
	z.Entities[e.GetUUID()] = e
	e.SetZone(z)
}

func (z *Zone) RemoveEntity(e model.Entity) {
	delete(z.Entities, e.GetUUID())
}

func (z *Zone) Update(dt float64) {
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
