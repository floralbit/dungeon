package network

import (
	"github.com/floralbit/dungeon/game/model"
	"github.com/google/uuid"
)

type serverEvent struct {
	Entity      *entityEvent        `json:"entity,omitempty"`
	Zone        *zoneEvent          `json:"zone,omitempty"`
	WorldObject *worldObjectEvent   `json:"world_object,omitempty"`
	Message     *serverMessageEvent `json:"message,omitempty"`
}

type entityEvent struct {
	UUID uuid.UUID `json:"uuid"`

	Update  model.Entity       `json:"update,omitemtpy"`
	Spawn   model.Entity       `json:"spawn,omitempty"`
	Despawn bool               `json:"despawn"`
	Die     bool               `json:"die"`
	Move    *entityMoveEvent   `json:"move,omitempty"`
	Chat    *entityChatEvent   `json:"chat,omitempty"`
	Attack  *entityAttackEvent `json:"attack,omitemtpy"`
}

type entityMoveEvent struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type entityChatEvent struct {
	Message string `json:"message"`
}

type entityAttackEvent struct {
	Target   uuid.UUID `json:"target"`
	TargetHP int       `json:"target_hp"`
	Hit      bool      `json:"hit"`
	Damage   int       `json:"damage"`
}

type zoneEvent struct {
	UUID uuid.UUID `json:"uuid"`

	Load model.Zone `json:"load,omitempty"`
}

type worldObjectEvent struct {
	UUID uuid.UUID `json:"uuid"`

	Spawn   model.WorldObject `json:"spawn,omitempty"`
	Despawn bool              `json:"despawn"`
}

type serverMessageEvent struct {
	Message string `json:"message"`
}

func newUpdateEvent(e model.Entity) serverEvent {
	return serverEvent{
		Entity: &entityEvent{
			UUID:   e.GetUUID(),
			Update: e,
		},
	}
}

func newSpawnEvent(e model.Entity) serverEvent {
	return serverEvent{
		Entity: &entityEvent{
			UUID:  e.GetUUID(),
			Spawn: e,
		},
	}
}

func newDespawnEvent(e model.Entity, becauseDeath bool) serverEvent {
	return serverEvent{
		Entity: &entityEvent{
			UUID:    e.GetUUID(),
			Despawn: true,
			Die:     becauseDeath,
		},
	}
}

func newMoveEvent(e model.Entity, x, y int) serverEvent {
	return serverEvent{
		Entity: &entityEvent{
			UUID: e.GetUUID(),
			Move: &entityMoveEvent{
				X: x,
				Y: y,
			},
		},
	}
}

func newChatEvent(e model.Entity, message string) serverEvent {
	return serverEvent{
		Entity: &entityEvent{
			UUID: e.GetUUID(),
			Chat: &entityChatEvent{
				Message: message,
			},
		},
	}
}

func newAttackEvent(e model.Entity, target uuid.UUID, hit bool, damage int, targetHP int) serverEvent {
	return serverEvent{
		Entity: &entityEvent{
			UUID: e.GetUUID(),
			Attack: &entityAttackEvent{
				Target:   target,
				Hit:      hit,
				Damage:   damage,
				TargetHP: targetHP,
			},
		},
	}
}

func newZoneLoadEvent(z model.Zone) serverEvent {
	return serverEvent{
		Zone: &zoneEvent{
			UUID: z.GetUUID(),
			Load: z,
		},
	}
}

func newWorldObjectSpawnEvent(o model.WorldObject) serverEvent {
	return serverEvent{
		WorldObject: &worldObjectEvent{
			UUID:  o.GetUUID(),
			Spawn: o,
		},
	}
}

func newWorldObjectDespawnEvent(o model.WorldObject) serverEvent {
	return serverEvent{
		WorldObject: &worldObjectEvent{
			UUID:    o.GetUUID(),
			Despawn: true,
		},
	}
}

func newServerMessageEvent(message string) serverEvent {
	return serverEvent{
		Message: &serverMessageEvent{
			Message: message,
		},
	}
}
