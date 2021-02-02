package game

import (
	"github.com/google/uuid"
)

var startingZoneUUID = uuid.MustParse("10f8b073-cbd7-46b7-a6e3-9cbdf68a933f")
var zones = loadZones()

type tile struct {
	ID    int    `json:"id"`
	Solid bool   `json:"solid"`
	Name  string `json:"name"`
}

type zone struct {
	UUID   uuid.UUID `json:"uuid"`
	Name   string    `json:"name"`
	Width  int       `json:"width"`
	Height int       `json:"height"`
	Tiles  []tile    `json:"tiles"`

	Entities     map[uuid.UUID]entity       `json:"entities"`
	WorldObjects map[uuid.UUID]*worldObject `json:"world_objects"`
}

func (z *zone) getTile(x, y int) *tile {
	if x < 0 || y < 0 || x >= z.Width || y >= z.Height {
		return nil
	}

	index := (z.Width * y) + x
	return &z.Tiles[index]
}

func (z *zone) getWorldObjects(x, y int) []*worldObject {
	objs := []*worldObject{}

	for _, obj := range z.WorldObjects {
		if obj.X == x && obj.Y == y {
			objs = append(objs, obj)
		}
	}

	return objs
}

func (z *zone) addEntity(e entity) {
	z.Entities[e.Data().UUID] = e
	e.Data().zone = z
	z.send(newSpawnEvent(e.Data()))
	e.Send(newZoneLoadEvent(z)) // send entity the zone data
}

func (z *zone) removeEntity(e entity) {
	delete(z.Entities, e.Data().UUID)
	z.send(newDespawnEvent(e.Data()))
}

func (z *zone) send(event serverEvent) {
	for _, e := range z.Entities {
		e.Send(event)
	}
}

func (z *zone) update(dt float64) {
	for _, e := range z.Entities {
		e.Update(dt)
	}
}
