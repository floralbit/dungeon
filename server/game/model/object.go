package model

import "github.com/google/uuid"

type WorldObject interface {
	GetUUID() uuid.UUID
}
