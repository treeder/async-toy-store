package models

import "encoding/json"

// Message wraps the payload to provide further metadata
// todo: convert this to a CloudEvent?
type Message struct {
	Channel      string          `json:"channel"`
	ReplyChannel string          `json:"reply_channel"`
	Payload      json.RawMessage `json:"payload"`
	// could we put an auth token here?
}
