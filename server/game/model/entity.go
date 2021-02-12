package model

import (
	"github.com/floralbit/dungeon/model"
	"github.com/google/uuid"
)

type Entity interface {
	GetUUID() uuid.UUID
	GetZone() Zone
	GetClient() *model.Client
}
