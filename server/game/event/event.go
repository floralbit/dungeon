package event

import "github.com/floralbit/dungeon/game/model"

type Event interface {
}

type JoinEvent struct {
	Entity model.Entity
}

type LeaveEvent struct {
	Entity model.Entity
}

type SpawnEvent struct {
	Entity model.Entity
}

type DespawnEvent struct {
	Entity model.Entity
}

type DieEvent struct {
	Entity model.Entity
}

type MoveEvent struct {
	Entity model.Entity
	X, Y   int
}

type ChatEvent struct {
	Entity  model.Entity
	Message string
}

type AttackEvent struct {
	Attacker, Target model.Entity
	Hit              bool
	Damage           int
	TargetHP         int
}

type HealEvent struct {
	Entity model.Entity
	Amount int
	Full   bool
}

type GainXPEvent struct {
	Entity    model.Entity
	LeveledUp bool
}
