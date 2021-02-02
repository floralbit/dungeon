package game

import (
	"github.com/google/uuid"
)

type entityType string

const (
	entityTypePlayer  = "player"
	entityTypeMonster = "monster"
)

type entity interface {
	Move(int, int)
	Spawn(uuid.UUID)
	Despawn()

	Send(serverEvent)

	Data() *entityData
}

type entityData struct {
	UUID uuid.UUID  `json:"uuid"`
	Name string     `json:"name"`
	Tile int        `json:"tile"` // representing tile
	Type entityType `json:"type"`

	Stats stats `json:"stats"`

	X int `json:"x"`
	Y int `json:"y"`

	zone *zone `json:"-"`
}

type stats struct {
	Level int `json:"level"`
	HP    int `json:"hp"`
	XP    int `json:"xp"'`
	AC    int `json:"ac"`

	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
}

func (e *entityData) Move(x, y int) {
	e.X = x
	e.Y = y

	e.zone.send(newMoveEvent(e, x, y))
}

func (e *entityData) Spawn(zoneUUID uuid.UUID) {
	z := zones[zoneUUID]
	z.addEntity(e)
}

func (e *entityData) Despawn() {
	e.zone.removeEntity(e)
}

func (e *entityData) Send(event serverEvent) {
	// NOP as default, players handle sends only
}

func (e *entityData) Data() *entityData {
	return e
}

func modifier(stat int) int {
	return (stat - 10) / 2
}
