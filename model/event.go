package model

// Event ...
type Event struct {
	// TODO: come up with better data model that works fine
	// for json serialization still lol, this fake union sucks
	ChatMessage *ChatMessage `json:"chatMessage,omitEmpty"`
}

// ChatMessage ...
type ChatMessage struct {
	Data string
}
