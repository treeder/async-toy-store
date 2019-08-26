package nats

import (
	"context"
	"encoding/json"
	"fmt"

	natsgo "github.com/nats-io/nats.go"
	"github.com/treeder/async-toy-store/brokers"
	"github.com/treeder/async-toy-store/models"
)

func Connect(ctx context.Context, urlStr string) (*Broker, error) {
	nc, err := natsgo.Connect(urlStr)
	if err != nil {
		return nil, err
	}
	b := &Broker{nc: nc}
	return b, nil
}

type Broker struct {
	nc *natsgo.Conn
}

func (b *Broker) Publish(ctx context.Context, channel string, msg *models.Message) error {
	// should this pass along the full wrapped Message or the Payload only?
	msgMarshalled, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return b.nc.Publish(channel, msgMarshalled)
}
func (b *Broker) Subscribe(ctx context.Context, channel string, handler brokers.Handler) (brokers.Subscription, error) {
	sub, err := b.nc.Subscribe(channel, func(m *natsgo.Msg) {
		msg, err := models.ParseMessage(m.Data)
		if err != nil {
			fmt.Printf("SUBSCRIBE ERROR: %v\n", err)
			// todo: be nice to have a way for user to deal with errors
			return
		}
		handler(msg)
	})
	return &Subscription{sub: sub}, err
}
func (b *Broker) Close() {
	b.nc.Close()
}

type Subscription struct {
	sub *natsgo.Subscription
}

func (s *Subscription) Unsubscribe() error {
	return s.sub.Unsubscribe()
}
