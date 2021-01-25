package model

// ServerEvent ...
type ServerEvent struct {
	Error *ServerErrorEvent `json:"error,omitempty"`

	Chat *ServerChatEvent `json:"chat,omitempty"`

	Zone *ServerZoneEvent `json:"zone,omitempty"`

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
	From string // username of joining player
}

// ServerLeaveEvent ...
type ServerLeaveEvent struct {
	From string // username of leaving player
}
