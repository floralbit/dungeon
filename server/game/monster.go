package game

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nickdavies/go-astar/astar"
	"log"
)

type monster struct {
	entityData

	moveSpeed    float64
	agroDistance float64

	moveTimer float64
}

type monsterType string

const (
	monsterTypeGoblin   = "goblin"
	monsterTypeSkeleton = "skeleton"
)

func newMonster(t monsterType) *monster {
	template, ok := monsterTemplates[t]
	if !ok {
		log.Fatal(errors.New(fmt.Sprintf("no monster found for type %s", t)))
	}

	m := &monster{
		entityData: entityData{
			UUID: uuid.New(),
			Name: template.Name,
			Tile: template.Tile,
			Type: entityTypeMonster,
			Stats: stats{
				Level: template.Level,

				Strength:     template.Strength,
				Dexterity:    template.Dexterity,
				Constitution: template.Constitution,
				Intelligence: template.Intelligence,
				Wisdom:       template.Wisdom,
				Charisma:     template.Charisma,
			},
		},
		moveSpeed:    template.MoveSpeed,
		agroDistance: template.AgroDistance,
	}

	// roll hp based on hd
	m.Stats.HP = 0
	for m.Stats.HP <= 0 {
		m.Stats.HP = roll{8, template.Level, modifier(m.Stats.Constitution)}.roll()
	}
	m.Stats.AC = 10 + modifier(m.Stats.Dexterity)

	return m
}

func (m *monster) Update(dt float64) {
	m.moveTimer += dt
	if m.moveTimer >= m.moveSpeed {
		m.moveTimer = 0
		m.move()
	}
}

func (m *monster) move() {
	// TODO: monster state machine w/ agro state on specific player
	for _, e := range m.zone.Entities {
		if e.Data().Type == entityTypePlayer {
			// just target the first player we see in the zone
			if dist(m.X, m.Y, e.Data().X, e.Data().Y) < m.agroDistance {
				// run a*
				a := astar.NewAStar(m.zone.Width, m.zone.Height)
				p2p := astar.NewPointToPoint()

				// avoid walls
				for x := 0; x < m.zone.Width; x++ {
					for y := 0; y < m.zone.Height; y++ {
						t := m.zone.getTile(x, y)
						if t.Solid {
							a.FillTile(astar.Point{Row: x, Col: y}, -1)
						}
					}
				}

				// avoid other monsters
				for _, otherE := range m.zone.Entities {
					if otherE.Data().Type == entityTypeMonster && otherE.Data().UUID != m.UUID {
						a.FillTile(astar.Point{Row: otherE.Data().X, Col: otherE.Data().Y}, -1)
					}
				}

				source := []astar.Point{{Row: m.X, Col: m.Y}}
				target := []astar.Point{{Row: e.Data().X, Col: e.Data().Y}}

				path := a.FindPath(p2p, source, target)
				if path != nil && path.Parent != nil {
					if path.Parent.Row == e.Data().X && path.Parent.Col == e.Data().Y {
						m.Attack(e)
					} else {
						m.Move(path.Parent.Row, path.Parent.Col) // move to first tile on path
					}
				}
				return
			}
		}
	}
}
