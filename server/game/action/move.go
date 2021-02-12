package action

import (
	"github.com/floralbit/dungeon/game/event"
	"github.com/floralbit/dungeon/game/model"
	"github.com/floralbit/dungeon/game/util"
)

const (
	maxValidMoveDist = 3
)

type MoveAction struct {
	Mover model.Entity
	X, Y  int
}

func (a *MoveAction) Execute() bool {
	eX, eY := a.Mover.GetPosition()
	if util.Dist(eX, eY, a.X, a.Y) > maxValidMoveDist {
		// can't move, probably a move race with zone warping
		event.NotifyObservers(event.MoveEvent{Entity: a.Mover, X: eX, Y: eY}) // tell them they're stationary
		return false
	}

	t := a.Mover.GetZone().GetTile(a.X, a.Y)
	if t == nil || t.Solid {
		// edge of map or solid, don't move
		event.NotifyObservers(event.MoveEvent{Entity: a.Mover, X: eX, Y: eY}) // tell them they're stationary
		return false
	}

	if a.Mover.GetType() == model.EntityTypePlayer {
		objs := a.Mover.GetZone().GetWorldObjects(a.X, a.Y)
		for _, obj := range objs {
			if obj.WarpTarget != nil {
				event.NotifyObservers(event.DespawnEvent{Entity: a.Mover})
				a.Mover.GetZone().RemoveEntity(a.Mover)
				obj.WarpTarget.Zone.AddEntity(a.Mover)
				a.Mover.SetPosition(obj.WarpTarget.X, obj.WarpTarget.Y)
				event.NotifyObservers(event.SpawnEvent{Entity: a.Mover})
				return true
			}
			if obj.HealZone != nil {
				if obj.HealZone.Full {
					a.Mover.Heal(a.Mover.GetStats().MaxHP)
					event.NotifyObservers(event.HealEvent{Entity: a.Mover, Amount: a.Mover.GetStats().MaxHP, Full: true})
					event.NotifyObservers(event.MoveEvent{Entity: a.Mover, X: eX, Y: eY}) // tell them they're stationary
					return true
				}
			}
		}
	}

	for _, otherE := range a.Mover.GetZone().GetEntities() {
		otherX, otherY := otherE.GetPosition()
		if a.Mover != otherE && otherX == a.X && otherY == a.Y {
			// someone is there, block the way
			event.NotifyObservers(event.MoveEvent{Entity: a.Mover, X: eX, Y: eY}) // tell them they're stationary
			return false
		}
	}

	a.Mover.SetPosition(a.X, a.Y)
	event.NotifyObservers(event.MoveEvent{Entity: a.Mover, X: a.X, Y: a.Y})
	return true // success
}
