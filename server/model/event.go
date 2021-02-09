package model

import "github.com/google/uuid"

// ClientEvent ...
type ClientEvent struct {
	Join   *ClientJoinEvent   `json:"join,omitempty"`
	Leave  *ClientLeaveEvent  `json:"leave,omitempty"`
	Chat   *ClientChatEvent   `json:"chat,omitempty"`
	Move   *ClientMoveEvent   `json:"move,omitempty"`
	Attack *ClientAttackEvent `json:"attack,omitempty"`

	Sender *Client `json:"-"`
}

// ClientJoinEvent ...
type ClientJoinEvent struct {
	Ok bool
}

// ClientLeaveEvent ...
type ClientLeaveEvent struct {
	Ok bool
}

// ClientChatEvent ...
type ClientChatEvent struct {
	Message string `json:"message"`
}

// ClientMoveEvent ...
type ClientMoveEvent struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// ClientAttackEvent ...
type ClientAttackEvent struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// ServerEvent ...
type ServerEvent struct {
	Connect *ServerConnectEvent `json:"connect,omitempty"`
	Error   *ServerErrorEvent   `json:"error,omitempty"`
}

// ServerConnectEvent
type ServerConnectEvent struct {
	UUID uuid.UUID `json:"uuid"`
}

// ServerErrorEvent ...
type ServerErrorEvent struct {
	Message string `json:"message"`
}
