package game

import "math/rand"

type roll struct {
	Sides, N, Plus int
}

func (r roll) roll() int {
	var total int
	for i := 0; i < r.N; i++ {
		total += rand.Intn(r.Sides)
	}
	return total + r.Plus
}
