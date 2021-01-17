package model

// ClientEvent ...
type ClientEvent struct {
	Chat     *ClientChatEvent     `json:"chat,omitempty"`
	Login    *ClientLoginEvent    `json:"login,omitempty"`
	Register *ClientRegisterEvent `json:"register,omitempty"`
	Leave    *ClientLeaveEvent    `json:"leave,omitempty"`

	Sender *Client `json:"-"`
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

// ClientLeaveEvent ...
type ClientLeaveEvent struct {
	Ok bool
}
