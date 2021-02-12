package game

import (
	"github.com/floralbit/dungeon/game/event"
	"github.com/floralbit/dungeon/game/model"
	"github.com/floralbit/dungeon/game/util"
	serverModel "github.com/floralbit/dungeon/model"
	"github.com/google/uuid"
)

type entityData struct {
	UUID uuid.UUID        `json:"uuid"`
	Name string           `json:"name"`
	Tile int              `json:"tile"` // representing tile
	Type model.EntityType `json:"type"`

	Stats model.Stats `json:"stats"`

	X int `json:"x"`
	Y int `json:"y"`

	EnergyThreshold int `json:"-"`
	Energy          int `json:"-"`

	queuedAction action     `json:"-"`
	zone         model.Zone `json:"-"`
}

func (e *entityData) GetUUID() uuid.UUID {
	return e.UUID
}

func (e *entityData) GetZone() model.Zone {
	return e.zone
}

func (e *entityData) SetZone(z model.Zone) {
	e.zone = z
}

func (e *entityData) GetType() model.EntityType {
	return e.Type
}

func (e *entityData) GetPosition() (int, int) {
	return e.X, e.Y
}

func (e *entityData) SetPosition(x, y int) {
	e.X = x
	e.Y = y
}

func (e *entityData) GetClient() *serverModel.Client {
	return nil
}

func (e *entityData) Act() model.Action {
	return nil // NOP
}

func (e *entityData) Tick() bool {
	e.Energy++
	if e.Energy >= e.EnergyThreshold {
		e.Energy = 0
		return true
	}
	return false
}

func (e *entityData) Die() {
	notifyObservers(event.DieEvent{Entity: e})
	e.zone.RemoveEntity(e)
}

func (e *entityData) GetStats() model.Stats {
	return e.Stats
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

func (e *entityData) RollToHit(targetAC int) bool {
	toHit := util.Roll{Sides: 20, N: 1, Plus: util.Modifier(e.Stats.Strength)}.Roll() // TODO: swap modifier based on weapon
	if toHit >= targetAC {
		return true
	}
	return false
}

func (e *entityData) RollDamage() int {
	damage := util.Roll{Sides: 3, N: 1, Plus: util.Modifier(e.Stats.Strength)}.Roll()
	if damage <= 0 {
		damage = 1 // minimum 1 dmg
	}
	return damage
}
