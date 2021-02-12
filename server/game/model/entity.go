package model

import (
	"github.com/floralbit/dungeon/model"
	"github.com/google/uuid"
)

type Entity interface {
	GetUUID() uuid.UUID
	GetType() EntityType

	GetZone() Zone
	SetZone(Zone)

	GetPosition() (int, int)
	SetPosition(int, int)

	GetStats() Stats

	Tick() bool // increments energy, returns if can act
	Act() Action

	RollToHit(int) bool
	RollDamage() int

	TakeDamage(int) bool
	Die()
	GainExp(int)
	Heal(int)

	GetClient() *model.Client
}

type EntityType string

const (
	EntityTypePlayer  = "player"
	EntityTypeMonster = "monster"
)

type Stats struct {
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
