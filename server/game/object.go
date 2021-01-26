package game

import (
	"github.com/google/uuid"
)

type worldObject struct {
	UUID uuid.UUID `json:"uuiud"`
	Name string    `json:"name"`
	Tile int       `json:"tile"` // representing tile
	X    int       `json:"int"`
	Y    int       `json:"int"`
}
