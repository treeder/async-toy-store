package brokers

import (
	"context"

	"github.com/treeder/async-toy-store/models"
)

// Broker is a common interface to underlying message brokers.
// todo: should we call this a broker or something better?
type Broker interface {
	// Publish ... change this to PublishMsg and have raw bytes as Publish
	Publish(ctx context.Context, topic string, message *models.Message) error
	Subscribe(ctx context.Context, topic string, handler Handler) (Subscription, error)
	Close()
}


type Subscription interface {
	Unsubscribe() error
}

type Message interface {
	Channel() string
	Payload() []byte
}

type Handler func(*models.Message)
