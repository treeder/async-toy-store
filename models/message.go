package models

import (
	"encoding/json"
	"fmt"
)

// Message wraps the payload to provide further metadata
// todo: convert this to a CloudEvent?
type Message struct {
	Channel      string          `json:"channel"`
	ReplyChannel string          `json:"reply_channel"`
	Payload      json.RawMessage `json:"payload"`
	// could we put an auth token here?
}

func ParseMessage(data []byte) (*Message, *Order, error) {
	msg := &Message{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		fmt.Println("error:", err)
		return nil, nil, err
	}
	fmt.Printf("MSG: %+v\n", msg)
	order, err := ParseOrder(msg.Payload)
	if err != nil {
		return msg, nil, err
	}
	return msg, order, nil
}

func ParseOrder(data []byte) (*Order, error) {
	order := &Order{}
	// can do a switch on msg.channel here to parse different object types
	err := json.Unmarshal(data, order)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}
	fmt.Printf("ORDER: %+v\n", order)
	return order, nil
}
