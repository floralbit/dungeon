package model

import "github.com/google/uuid"

type WorldObject struct {
	UUID uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
	Tile int       `json:"tile"` // representing tile
	X    int       `json:"x"`
	Y    int       `json:"y"`

	// special features
	Type       WorldObjectType `json:"type"`
	WarpTarget *WarpTarget     `json:"warp_target,omitemtpy"`
	HealZone   *HealZone       `json:"heal_zone,omitempty"`
}

type WorldObjectType string

const (
	WorldObjectTypePlayerSpawn = "playerSpawn"
	WorldObjectTypePortal      = "portal"
	WorldObjectTypeHealing     = "healing"
)

type WarpTarget struct {
	ZoneUUID uuid.UUID `json:"zone_uuid"`
	Zone     Zone      `json:"-"`
	X        int       `json:"x"`
	Y        int       `json:"y"`
}

type HealZone struct {
	Full bool
}
