package model

// ServerEvent ...
type ServerEvent struct {
	Error *ServerErrorEvent `json:"error,omitempty"`

	Zone *ServerZoneEvent `json:"zone,omitempty"`
	Chat *ServerChatEvent `json:"chat,omitempty"`
	Move *ServerMoveEvent `json:"move,omitempty"`

	Join  *ServerJoinEvent  `json:"join,omitempty"`
	Leave *ServerLeaveEvent `json:"leave,omitempty"`
}

// ServerErrorEvent ...
type ServerErrorEvent struct {
	Message string
}

// ServerChatEvent ...
type ServerChatEvent struct {
	Message string
	From    string // username of sender
}

// ServerZoneEvent ...
type ServerZoneEvent struct {
	Data interface{}
}

// ServerJoinEvent ...
type ServerJoinEvent struct {
	Data interface{} // data of joining player
	From string      // username of joining player
	You  bool        // set for joining player to know it's them
}

// ServerLeaveEvent ...
type ServerLeaveEvent struct {
	From string // username of leaving player
}

// ServerMoveEvent ...
type ServerMoveEvent struct {
	X, Y int
	From string // username of moving player
	You  bool
}
