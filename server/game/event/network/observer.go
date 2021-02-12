package network

import "github.com/floralbit/dungeon/game/event"

const MOTD = "Welcome to Vault of Splendor adventurer! Stay safe."

type networkObserver struct {
}

func NewObserver() event.Observer {
	return &networkObserver{}
}

func (o *networkObserver) Notify(e event.Event) {
	switch v := e.(type) {
	case event.JoinEvent:
		if c := v.Entity.GetClient(); c != nil {
			c.In <- newServerMessageEvent(MOTD)
		}

	case event.LeaveEvent:
		break

	case event.SpawnEvent:
		for _, ent := range v.Entity.GetZone().GetEntities() {
			if c := ent.GetClient(); c != nil {
				c.In <- newSpawnEvent(v.Entity)
			}
		}
		if c := v.Entity.GetClient(); c != nil {
			c.In <- newZoneLoadEvent(v.Entity.GetZone())
		}

	case event.DespawnEvent:
		for _, ent := range v.Entity.GetZone().GetEntities() {
			if c := ent.GetClient(); c != nil {
				c.In <- newDespawnEvent(v.Entity, false)
			}
		}

	case event.DieEvent:
		for _, ent := range v.Entity.GetZone().GetEntities() {
			if c := ent.GetClient(); c != nil {
				c.In <- newDespawnEvent(v.Entity, true)
			}
		}
		if c := v.Entity.GetClient(); c != nil {
			c.In <- newServerMessageEvent("You died.")
			c.In <- newServerMessageEvent("Your soul enters a new body. You are reborn.")
		}

	case event.MoveEvent:
		for _, ent := range v.Entity.GetZone().GetEntities() {
			if c := ent.GetClient(); c != nil {
				c.In <- newMoveEvent(v.Entity, v.X, v.Y)
			}
		}

	case event.ChatEvent:
		for _, ent := range v.Entity.GetZone().GetEntities() {
			if c := ent.GetClient(); c != nil {
				c.In <- newChatEvent(v.Entity, v.Message)
			}
		}

	case event.AttackEvent:
		for _, ent := range v.Attacker.GetZone().GetEntities() {
			if c := ent.GetClient(); c != nil {
				c.In <- newAttackEvent(v.Attacker, v.Target.GetUUID(), v.Hit, v.Damage, v.TargetHP)
			}
		}

	case event.HealEvent:
		if c := v.Entity.GetClient(); c != nil {
			c.In <- newUpdateEvent(v.Entity)
			c.In <- newServerMessageEvent("You pray to your gods and are fully healed in their light.")
		}

	case event.GainXPEvent:
		if c := v.Entity.GetClient(); c != nil {
			c.In <- newUpdateEvent(v.Entity)
			if v.LeveledUp {
				c.In <- newServerMessageEvent("You leveled up! You have a newfound strength coursing through your veins.")
			}
		}

	}
}
