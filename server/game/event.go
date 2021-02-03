package game

import "github.com/google/uuid"

type serverEvent struct {
	Entity      *entityEvent        `json:"entity,omitempty"`
	Zone        *zoneEvent          `json:"zone,omitempty"`
	WorldObject *worldObjectEvent   `json:"world_object,omitempty"`
	Message     *serverMessageEvent `json:"message,omitempty"`
}

type entityEvent struct {
	UUID uuid.UUID `json:"uuid"`

	Spawn   *entityData        `json:"spawn,omitempty"`
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

	Load *zone `json:"load,omitempty"`
}

type worldObjectEvent struct {
	UUID uuid.UUID `json:"uuid"`

	Spawn   *worldObject `json:"spawn,omitempty"`
	Despawn bool         `json:"despawn"`
}

type serverMessageEvent struct {
	Message string `json:"message"`
}

func newSpawnEvent(e *entityData) serverEvent {
	return serverEvent{
		Entity: &entityEvent{
			UUID:  e.UUID,
			Spawn: e,
		},
	}
}

func newDespawnEvent(e *entityData, becauseDeath bool) serverEvent {
	return serverEvent{
		Entity: &entityEvent{
			UUID:    e.UUID,
			Despawn: true,
			Die:     becauseDeath,
		},
	}
}

func newMoveEvent(e *entityData, x, y int) serverEvent {
	return serverEvent{
		Entity: &entityEvent{
			UUID: e.UUID,
			Move: &entityMoveEvent{
				X: x,
				Y: y,
			},
		},
	}
}

func newChatEvent(e *entityData, message string) serverEvent {
	return serverEvent{
		Entity: &entityEvent{
			UUID: e.UUID,
			Chat: &entityChatEvent{
				Message: message,
			},
		},
	}
}

func newAttackEvent(e *entityData, target uuid.UUID, hit bool, damage int, targetHP int) serverEvent {
	return serverEvent{
		Entity: &entityEvent{
			UUID: e.UUID,
			Attack: &entityAttackEvent{
				Target:   target,
				Hit:      hit,
				Damage:   damage,
				TargetHP: targetHP,
			},
		},
	}
}

func newZoneLoadEvent(z *zone) serverEvent {
	return serverEvent{
		Zone: &zoneEvent{
			UUID: z.UUID,
			Load: z,
		},
	}
}

func newWorldObjectSpawnEvent(o *worldObject) serverEvent {
	return serverEvent{
		WorldObject: &worldObjectEvent{
			UUID:  o.UUID,
			Spawn: o,
		},
	}
}

func newWorldObjectDespawnEvent(o *worldObject) serverEvent {
	return serverEvent{
		WorldObject: &worldObjectEvent{
			UUID:    o.UUID,
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
