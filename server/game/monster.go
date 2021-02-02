package game

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
)

type monster struct {
	entityData
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
