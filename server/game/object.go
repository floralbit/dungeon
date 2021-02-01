package game

import (
	"github.com/google/uuid"
)

type worldObject struct {
	UUID uuid.UUID `json:"uuiud"`
	Name string    `json:"name"`
	Tile int       `json:"tile"` // representing tile
	X    int       `json:"x"`
	Y    int       `json:"y"`

	// special features
	Type       worldObjectType `json:"type"`
	WarpTarget *warpTarget     `json:"warp_target,omitemtpy"`
}

type worldObjectType string

const (
	worldObjectTypePlayerSpawn = "playerSpawn"
)

type warpTarget struct {
	ZoneUUID uuid.UUID `json:"zone_uuid"`
	X        int       `json:"x"`
	Y        int       `json:"y"`
}
