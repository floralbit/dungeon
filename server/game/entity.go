package game

import (
	"github.com/google/uuid"
	"math"
)

type entityType string

const (
	entityTypePlayer  = "player"
	entityTypeMonster = "monster"
)

const (
	xpLevelFactor = 500
)

type entity interface {
	Update(dt float64)
	Send(serverEvent)

	Spawn(uuid.UUID)
	Despawn(bool)
	Die()

	GainExp(int)
	TakeDamage(int) bool
	Heal(int)
	HealFull()

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

func (e *entityData) Update(dt float64) {
	// NOP as default
}

func (e *entityData) Send(event serverEvent) {
	// NOP as default, players handle sends only
}

func (e *entityData) Spawn(zoneUUID uuid.UUID) {
	z := zones[zoneUUID]
	z.addEntity(e)
}

func (e *entityData) Despawn(becauseDeath bool) {
	e.zone.removeEntity(e, becauseDeath)
}

func (e *entityData) Die() {
	e.Despawn(true)
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
	nextLevelXP := xpForLevel(e.Stats.Level)
	for e.Stats.XP >= nextLevelXP {
		e.Stats.Level += 1
		e.Stats.MaxHP += roll{8, 1, modifier(e.Stats.Constitution)}.roll()
		e.Stats.HP = e.Stats.MaxHP
		nextLevelXP = xpForLevel(e.Stats.Level)
	}
	e.Stats.XPToNextLevel = nextLevelXP
}

func (e *entityData) Heal(amount int) {
	e.Stats.HP += amount
	if e.Stats.HP > e.Stats.MaxHP {
		e.Stats.HP = e.Stats.MaxHP
	}
}

func (e *entityData) HealFull() {
	e.Stats.HP = e.Stats.MaxHP
}

func (e *entityData) Data() *entityData {
	return e
}

func (e *entityData) rollToHit(targetAC int) bool {
	toHit := roll{Sides: 20, N: 1, Plus: modifier(e.Stats.Strength)}.roll() // TODO: swap modifier based on weapon
	if toHit >= targetAC {
		return true
	}
	return false
}

func (e *entityData) rollDamage() int {
	damage := roll{Sides: 3, N: 1, Plus: modifier(e.Stats.Strength)}.roll()
	if damage <= 0 {
		damage = 1 // minimum 1 dmg
	}
	return damage
}

func modifier(stat int) int {
	return (stat - 10) / 2
}

func worthXP(level int) int {
	return level * 100
}

func xpForLevel(level int) int {
	return int(xpLevelFactor*math.Pow(float64(level), 2) - float64(xpLevelFactor*level))
}
