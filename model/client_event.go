package model

// ClientEvent ...
type ClientEvent struct {
	Chat  *ClientChatEvent  `json:"chat,omitempty"`
	Login *ClientLoginEvent `json:"login,omitempty"`
}

// ClientChatEvent ...
type ClientChatEvent struct {
	Message string
}

// ClientLoginEvent ...
type ClientLoginEvent struct {
	Username string
	Password string
}

// ClientRegisterEvent ...
type ClientRegisterEvent struct {
	Username string
	Password string
}
