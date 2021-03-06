package entity

import (
	"errors"
	"fmt"
	"github.com/floralbit/dungeon/game/action"
	"github.com/floralbit/dungeon/game/data"
	"github.com/floralbit/dungeon/game/model"
	"github.com/floralbit/dungeon/game/util"
	"github.com/google/uuid"
	"github.com/nickdavies/go-astar/astar"
	"log"
)

type Monster struct {
	entityData

	agroDistance float64
}

type monsterType string

const (
	MonsterTypeGoblin   = "goblin"
	MonsterTypeSkeleton = "skeleton"
)

func NewMonster(t monsterType) *Monster {
	template, ok := data.MonsterTemplates[string(t)]
	if !ok {
		log.Fatal(errors.New(fmt.Sprintf("no Monster found for type %s", t)))
	}

	m := &Monster{
		entityData: entityData{
			UUID: uuid.New(),
			Name: template.Name,
			Tile: template.Tile,
			Type: model.EntityTypeMonster,
			Stats: model.Stats{
				Level: template.Level,

				Strength:     template.Strength,
				Dexterity:    template.Dexterity,
				Constitution: template.Constitution,
				Intelligence: template.Intelligence,
				Wisdom:       template.Wisdom,
				Charisma:     template.Charisma,
			},
			EnergyThreshold: template.EnergyThreshold,
		},
		agroDistance: template.AgroDistance,
	}

	// roll hp based on hd
	m.Stats.HP = 0
	for m.Stats.HP <= 0 {
		m.Stats.HP = util.Roll{8, template.Level, util.Modifier(m.Stats.Constitution)}.Roll()
	}
	m.Stats.MaxHP = m.Stats.HP
	m.Stats.AC = 10 + util.Modifier(m.Stats.Dexterity)

	return m
}

func (m *Monster) Act() model.Action {
	// TODO: Monster state machine w/ agro state on specific player
	for _, e := range m.zone.GetEntities() {
		if e.GetType() == model.EntityTypePlayer {
			// just target the first player we see in the zone
			eX, eY := e.GetPosition()
			if util.Dist(m.X, m.Y, eX, eY) < m.agroDistance {
				// run a*
				w, h := m.zone.GetDimensions()
				a := astar.NewAStar(w, h)
				p2p := astar.NewPointToPoint()

				// avoid walls
				for x := 0; x < w; x++ {
					for y := 0; y < h; y++ {
						t := m.zone.GetTile(x, y)
						if t.Solid {
							a.FillTile(astar.Point{Row: x, Col: y}, -1)
						}
					}
				}

				// avoid other monsters by looking at their current pos or planned movement
				for _, otherE := range m.zone.GetEntities() {
					if otherE.GetType() == model.EntityTypeMonster && otherE.GetUUID() != m.UUID {
						otherM := otherE.(*Monster)
						if otherMove, ok := otherM.QueuedAction.(*action.MoveAction); ok {
							a.FillTile(astar.Point{Row: otherMove.X, Col: otherMove.Y}, 3)
						} else {
							a.FillTile(astar.Point{Row: otherM.X, Col: otherM.Y}, 3)
						}
					}
				}

				source := []astar.Point{{Row: m.X, Col: m.Y}}
				target := []astar.Point{{Row: eX, Col: eY}}

				path := a.FindPath(p2p, source, target)
				if path != nil && path.Parent != nil {
					if path.Parent.Row == eX && path.Parent.Col == eY {
						m.QueuedAction = &action.LightAttackAction{
							Attacker: m,
							X:        path.Parent.Row,
							Y:        path.Parent.Col,
						}
						return m.QueuedAction
					} else {
						m.QueuedAction = &action.MoveAction{
							Mover: m,
							X:     path.Parent.Row,
							Y:     path.Parent.Col,
						}
						return m.QueuedAction
					}
				}
			}
		}
	}
	m.QueuedAction = nil
	return nil
}
