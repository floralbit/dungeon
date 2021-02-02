package game

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
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
				dx := e.Data().X - m.X
				dy := e.Data().Y - m.Y
				var moveX, moveY int
				if dx > 0 {
					moveX = 1
				}
				if dx < 0 {
					moveX = -1
				}
				if dy > 0 {
					moveY = 1
				}
				if dy < 0 {
					moveY = -1
				}
				m.Move(m.X+moveX, m.Y+moveY)
				return
			}
		}
	}
}
