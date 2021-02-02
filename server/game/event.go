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

	Spawn   *entityData      `json:"spawn,omitempty"`
	Despawn bool             `json:"despawn"`
	Die     bool             `json:"die"`
	Move    *entityMoveEvent `json:"move,omitempty"`
	Chat    *entityChatEvent `json:"chat,omitempty"`
}

type entityMoveEvent struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type entityChatEvent struct {
	Message string `json:"message"`
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

func newDespawnEvent(e *entityData) serverEvent {
	return serverEvent{
		Entity: &entityEvent{
			UUID:    e.UUID,
			Despawn: true,
		},
	}
}

func newDieEvent(e *entityData) serverEvent {
	return serverEvent{
		Entity: &entityEvent{
			UUID: e.UUID,
			Die:  true,
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
