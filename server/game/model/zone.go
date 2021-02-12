package model

import "github.com/google/uuid"

type Zone interface {
	GetUUID() uuid.UUID
	GetEntities() []Entity
}
