package game

import "github.com/floralbit/dungeon/game/util"

type action interface {
	// Execute preforms the action
	Execute() bool
}

type lightAttackAction struct {
	Attacker entity
	X, Y     int
}

func (a *lightAttackAction) Execute() bool {
	var target entity
	for _, e := range a.Attacker.Data().zone.Entities {
		if e != a.Attacker && e.Data().X == a.X && e.Data().Y == a.Y {
			target = e
		}
	}
	if target == nil {
		return false // no target at location
	}

	var damage int
	var hit bool
	if a.Attacker.Data().rollToHit(target.Data().Stats.AC) {
		hit = true
		damage = a.Attacker.Data().rollDamage()
	}

	// resolve damage
	wouldDie := target.TakeDamage(damage)
	a.Attacker.Data().zone.send(newAttackEvent(a.Attacker.Data(), target.Data().UUID, hit, damage, target.Data().Stats.HP))

	if wouldDie {
		target.Die()
		a.Attacker.GainExp(util.WorthXP(target.Data().Stats.Level))
	}

	return true // success
}

type heavyAttackAction struct {
	Attacker, Target entity
}

func (a *heavyAttackAction) Execute() bool {
	return true // success
}

type moveAction struct {
	Mover entity
	X, Y  int
}

func (a *moveAction) Execute() bool {
	e := a.Mover.Data()
	if util.Dist(e.X, e.Y, a.X, a.Y) > maxValidMoveDist {
		// can't move, probably a move race with zone warping
		a.Mover.Send(newMoveEvent(e, e.X, e.Y)) // tell them they're stationary
		return false
	}

	t := e.zone.getTile(a.X, a.Y)
	if t == nil || t.Solid {
		// edge of map or solid, don't move
		a.Mover.Send(newMoveEvent(e, e.X, e.Y)) // tell them they're stationary
		return false
	}

	if p, ok := a.Mover.(*player); ok {
		objs := p.zone.getWorldObjects(a.X, a.Y)
		for _, obj := range objs {
			if obj.WarpTarget != nil {
				p.zone.removeEntity(p, false)
				zones[obj.WarpTarget.ZoneUUID].addEntity(p)
				p.X = obj.WarpTarget.X
				p.Y = obj.WarpTarget.Y
				return true
			}
			if obj.HealZone != nil {
				if obj.HealZone.Full {
					p.Heal(p.Stats.MaxHP)
					p.Send(newUpdateEvent(e))
					p.Send(newMoveEvent(e, p.X, p.Y)) // tell them they're stationary
					p.Send(newServerMessageEvent("You pray to your gods and are fully healed in their light."))
					return true
				}
			}
		}
	}

	for _, otherE := range a.Mover.Data().zone.Entities {
		if a.Mover != otherE && otherE.Data().X == a.X && otherE.Data().Y == a.Y {
			// someone is there, block the way
			a.Mover.Send(newMoveEvent(e, e.X, e.Y)) // tell them they're stationary
			return false
		}
	}

	e.X = a.X
	e.Y = a.Y

	e.zone.send(newMoveEvent(e, a.X, a.Y))
	return true // success
}
