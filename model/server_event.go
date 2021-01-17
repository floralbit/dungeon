package model

// ServerEvent ...
type ServerEvent struct {
	Error *ServerErrorEvent `json:"error,omitempty"`

	Chat     *ServerChatEvent     `json:"chat,omitempty"`
	Login    *ServerLoginEvent    `json:"login,omitempty"`
	Register *ServerRegisterEvent `json:"register,omitempty"`
	Leave    *ServerLeaveEvent    `json:"leave,omitempty"`
}

// ServerErrorEvent ...
type ServerErrorEvent struct {
	Message string
}

// ServerChatEvent ...
type ServerChatEvent struct {
	Message string
	From    string // UUID of sender
}

// ServerLoginEvent ...
type ServerLoginEvent struct {
	Ok bool
	// todo: fill in account data
}

// ServerRegisterEvent ...
type ServerRegisterEvent struct {
	Ok bool
}

// ServerLeaveEvent ...
type ServerLeaveEvent struct {
	From string // UUID of leaving player
}
