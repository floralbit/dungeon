package gen

import "math/rand"

type Object struct {
	Type ObjectType
}

type ObjectType string

const (
	ObjectTypeMonsterSlot = "monster_slot"

	monsterSlotLiklihood = 0.005 // .5%
)

func (l *Level) placeMonsters() {
	for x := 0; x < l.Width; x++ {
		for y := 0; y < l.Height; y++ {
			if !l.freeSpace(x, y) {
				continue
			}
			if rand.Float32() < monsterSlotLiklihood {
				l.Objects[x][y] = &Object{Type: ObjectTypeMonsterSlot}
			}
		}
	}
}
