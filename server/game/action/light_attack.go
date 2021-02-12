package action

import (
	"github.com/floralbit/dungeon/game/event"
	"github.com/floralbit/dungeon/game/model"
	"github.com/floralbit/dungeon/game/util"
)

type LightAttackAction struct {
	Attacker model.Entity
	X, Y     int
}

func (a *LightAttackAction) Execute() bool {
	var target model.Entity
	for _, e := range a.Attacker.GetZone().GetEntities() {
		eX, eY := e.GetPosition()
		if e != a.Attacker && eX == a.X && eY == a.Y {
			target = e
		}
	}
	if target == nil {
		if a.Attacker.GetType() == model.EntityTypePlayer {
			m := &MoveAction{Mover: a.Attacker, X: a.X, Y: a.Y}
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
	event.NotifyObservers(event.AttackEvent{Attacker: a.Attacker, Target: target, Hit: hit, Damage: damage, TargetHP: target.GetStats().HP})

	if wouldDie {
		target.Die()
		a.Attacker.GainExp(util.WorthXP(target.GetStats().Level))
	}

	return true // success
}
