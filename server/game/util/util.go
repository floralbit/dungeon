package util

import "math"

const (
	xpLevelFactor = 500
)

func Dist(x1, y1, x2, y2 int) float64 {
	return math.Sqrt(math.Pow(float64(x1-x2), 2) + math.Pow(float64(y1-y2), 2))
}

func Modifier(stat int) int {
	return (stat - 10) / 2
}

func WorthXP(level int) int {
	return level * 100
}

func XPForLevel(level int) int {
	return int(xpLevelFactor*math.Pow(float64(level), 2) - float64(xpLevelFactor*level))
}
