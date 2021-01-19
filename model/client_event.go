package model

// ClientEvent ...
type ClientEvent struct {
	Chat *ClientChatEvent `json:"chat,omitempty"`

	Join  *ClientJoinEvent  `json:"join,omitempty"`
	Leave *ClientLeaveEvent `json:"leave,omitempty"`

	Sender *Client `json:"-"`
}

// ClientChatEvent ...
type ClientChatEvent struct {
	Message string
}

// ClientJoinEvent ...
type ClientJoinEvent struct {
	Ok bool
}

// ClientLeaveEvent ...
type ClientLeaveEvent struct {
	Ok bool
}
