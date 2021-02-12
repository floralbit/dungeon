package game

import (
	"github.com/floralbit/dungeon/game/event"
	"github.com/floralbit/dungeon/game/model"
	"github.com/floralbit/dungeon/game/util"
	serverModel "github.com/floralbit/dungeon/model"
	"github.com/google/uuid"
)

type entityType string

const (
	entityTypePlayer  = "player"
	entityTypeMonster = "monster"
)

type entity interface {
	model.Entity

	Act() action

	Die()

	GainExp(int)
	TakeDamage(int) bool
	Heal(int)

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

	EnergyThreshold int `json:"-"`
	Energy          int `json:"-"`

	queuedAction action `json:"-"`
	zone         *zone  `json:"-"`
}

type stats struct {
	Level         int `json:"level"`
	MaxHP         int `json:"max_hp"`
	HP            int `json:"hp"`
	XP            int `json:"xp"`
	AC            int `json:"ac"`
	XPToNextLevel int `json:"xp_to_next_level"`

	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
}

func (e *entityData) GetUUID() uuid.UUID {
	return e.UUID
}

func (e *entityData) GetZone() model.Zone {
	return e.zone
}

func (e *entityData) GetClient() *serverModel.Client {
	return nil
}

func (e *entityData) Act() action {
	return nil // NOP
}

func (e *entityData) Die() {
	notifyObservers(event.DieEvent{Entity: e})
	e.zone.removeEntity(e)
}

// TakeDamage returns if they would die so XP can be dished out
func (e *entityData) TakeDamage(damage int) bool {
	e.Stats.HP -= damage
	if e.Stats.HP <= 0 {
		return true
	}
	return false
}

func (e *entityData) GainExp(xp int) {
	e.Stats.XP += xp
	nextLevelXP := util.XPForLevel(e.Stats.Level)
	for e.Stats.XP >= nextLevelXP {
		e.Stats.Level += 1
		e.Stats.MaxHP += util.Roll{8, 1, util.Modifier(e.Stats.Constitution)}.Roll()
		e.Stats.HP = e.Stats.MaxHP
		nextLevelXP = util.XPForLevel(e.Stats.Level)
	}
	e.Stats.XPToNextLevel = nextLevelXP
}

func (e *entityData) Heal(amount int) {
	e.Stats.HP += amount
	if e.Stats.HP > e.Stats.MaxHP {
		e.Stats.HP = e.Stats.MaxHP
	}
}

func (e *entityData) Data() *entityData {
	return e
}

func (e *entityData) rollToHit(targetAC int) bool {
	toHit := util.Roll{Sides: 20, N: 1, Plus: util.Modifier(e.Stats.Strength)}.Roll() // TODO: swap modifier based on weapon
	if toHit >= targetAC {
		return true
	}
	return false
}

func (e *entityData) rollDamage() int {
	damage := util.Roll{Sides: 3, N: 1, Plus: util.Modifier(e.Stats.Strength)}.Roll()
	if damage <= 0 {
		damage = 1 // minimum 1 dmg
	}
	return damage
}
