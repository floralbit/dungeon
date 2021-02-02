package game

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nickdavies/go-astar/astar"
	"log"
)

const (
	monsterMoveTime = 1.0 // in s
	monsterAgroDist = 6   // eucledian dist
)

type monster struct {
	entityData

	moveTimer float64
}

type monsterType string

const (
	monsterTypeGoblin = "goblin"
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
				Level:        template.HD,
				AC:           template.AC,
				Strength:     template.Strength,
				Dexterity:    template.Dexterity,
				Constitution: template.Constitution,
				Intelligence: template.Intelligence,
				Wisdom:       template.Wisdom,
				Charisma:     template.Charisma,
			},
		},
	}

	// roll hp based on hd
	m.Stats.HP = 0
	for m.Stats.HP <= 0 {
		m.Stats.HP = roll{8, template.HD, 0}.roll()
	}

	return m
}

func (m *monster) Update(dt float64) {
	m.moveTimer += dt
	if m.moveTimer >= monsterMoveTime {
		m.moveTimer = 0
		m.move()
	}
}

func (m *monster) move() {
	// TODO: monster state machine w/ agro state on specific player
	for _, e := range m.zone.Entities {
		if e.Data().Type == entityTypePlayer {
			if dist(m.X, m.Y, e.Data().X, e.Data().Y) < monsterAgroDist {
				a := astar.NewAStar(m.zone.Width, m.zone.Height)
				p2p := astar.NewPointToPoint()
				for x := 0; x < m.zone.Width; x++ {
					for y := 0; y < m.zone.Height; y++ {
						t := m.zone.getTile(x, y)
						if t.Solid {
							a.FillTile(astar.Point{Row: x, Col: y}, -1)
						}
					}
				}

				source := []astar.Point{{Row: m.X, Col: m.Y}}
				target := []astar.Point{{Row: e.Data().X, Col: e.Data().Y}}

				path := a.FindPath(p2p, source, target)
				if path != nil && path.Parent != nil {
					m.Move(path.Parent.Row, path.Parent.Col)
				}
				return
			}
		}
	}
}
