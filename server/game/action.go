package game

import (
	"github.com/floralbit/dungeon/game/event"
	"github.com/floralbit/dungeon/game/model"
	"github.com/floralbit/dungeon/game/util"
)

type action interface {
	// Execute preforms the action
	Execute() bool
}

type lightAttackAction struct {
	Attacker model.Entity
	X, Y     int
}

func (a *lightAttackAction) Execute() bool {
	var target model.Entity
	for _, e := range a.Attacker.GetZone().GetEntities() {
		eX, eY := e.GetPosition()
		if e != a.Attacker && eX == a.X && eY == a.Y {
			target = e
		}
	}
	if target == nil {
		if a.Attacker.GetType() == model.EntityTypePlayer {
			m := &moveAction{Mover: a.Attacker, X: a.X, Y: a.Y}
			return m.Execute() // for player case, move them if no target
		}
		return false // no target at location
	}

	var damage int
	var hit bool
	if a.Attacker.RollToHit(target.GetStats().AC) {
		hit = true
		damage = a.Attacker.RollDamage()
	}

	// resolve damage
	wouldDie := target.TakeDamage(damage)
	notifyObservers(event.AttackEvent{Attacker: a.Attacker, Target: target, Hit: hit, Damage: damage, TargetHP: target.GetStats().HP})

	if wouldDie {
		target.Die()
		a.Attacker.GainExp(util.WorthXP(target.GetStats().Level))
	}

	return true // success
}

type heavyAttackAction struct {
	Attacker, Target model.Entity
}

func (a *heavyAttackAction) Execute() bool {
	return true // success
}

type moveAction struct {
	Mover model.Entity
	X, Y  int
}

func (a *moveAction) Execute() bool {
	eX, eY := a.Mover.GetPosition()
	if util.Dist(eX, eY, a.X, a.Y) > maxValidMoveDist {
		// can't move, probably a move race with zone warping
		notifyObservers(event.MoveEvent{Entity: a.Mover, X: eX, Y: eY}) // tell them they're stationary
		return false
	}

	t := a.Mover.GetZone().GetTile(a.X, a.Y)
	if t == nil || t.Solid {
		// edge of map or solid, don't move
		notifyObservers(event.MoveEvent{Entity: a.Mover, X: eX, Y: eY}) // tell them they're stationary
		return false
	}

	if a.Mover.GetType() == model.EntityTypePlayer {
		objs := a.Mover.GetZone().GetWorldObjects(a.X, a.Y)
		for _, obj := range objs {
			if obj.WarpTarget != nil {
				notifyObservers(event.DespawnEvent{Entity: a.Mover})
				a.Mover.GetZone().RemoveEntity(a.Mover)
				zones[obj.WarpTarget.ZoneUUID].AddEntity(a.Mover)
				a.Mover.SetPosition(obj.WarpTarget.X, obj.WarpTarget.Y)
				notifyObservers(event.SpawnEvent{Entity: a.Mover})
				return true
			}
			if obj.HealZone != nil {
				if obj.HealZone.Full {
					a.Mover.Heal(a.Mover.GetStats().MaxHP)
					notifyObservers(event.HealEvent{Entity: a.Mover, Amount: a.Mover.GetStats().MaxHP, Full: true})
					notifyObservers(event.MoveEvent{Entity: a.Mover, X: eX, Y: eY}) // tell them they're stationary
					return true
				}
			}
		}
	}

	for _, otherE := range a.Mover.GetZone().GetEntities() {
		otherX, otherY := otherE.GetPosition()
		if a.Mover != otherE && otherX == a.X && otherY == a.Y {
			// someone is there, block the way
			notifyObservers(event.MoveEvent{Entity: a.Mover, X: eX, Y: eY}) // tell them they're stationary
			return false
		}
	}

	a.Mover.SetPosition(a.X, a.Y)
	notifyObservers(event.MoveEvent{Entity: a.Mover, X: a.X, Y: a.Y})
	return true // success
}
