package model

// ServerEvent ...
type ServerEvent struct {
	Error *ServerErrorEvent `json:"error,omitempty"`

	Chat     *ServerChatEvent     `json:"chat,omitempty"`
	Login    *ServerLoginEvent    `json:"login,omitempty"`
	Register *ServerRegisterEvent `json:"register,omitempty"`
}

// ServerErrorEvent ...
type ServerErrorEvent struct {
	Message string
}

// ServerChatEvent ...
type ServerChatEvent struct {
	Message string
	From    string
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
