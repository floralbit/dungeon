package util

import "math/rand"

type Roll struct {
	Sides, N, Plus int
}

func (r Roll) Roll() int {
	var total int
	for i := 0; i < r.N; i++ {
		total += rand.Intn(r.Sides) + 1
	}
	return total + r.Plus
}
