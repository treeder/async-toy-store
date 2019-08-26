package mqtt

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/treeder/async-toy-store/brokers"
	"github.com/treeder/async-toy-store/models"
)

func Connect(ctx context.Context, urlStr string) (*Broker, error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	// this also supports mqtt+ws URL scheme
	s := url.Scheme
	sp := strings.Split(s, "+")
	if len(sp) > 1 {
		url.Scheme = sp[1]
		urlStr = url.String()
	}
	mqttClient := mqtt.NewClient(mqtt.NewClientOptions().AddBroker(urlStr))
	token := mqttClient.Connect()
	token.Wait()
	if token.Error() != nil {
		return nil, token.Error()
	}

	b := &Broker{client: mqttClient}
	return b, nil
}

type Broker struct {
	client mqtt.Client
}

func (b *Broker) Publish(ctx context.Context, channel string, msg *models.Message) error {
	msgMarshalled, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	token2 := b.client.Publish(channel, 1, false, msgMarshalled)
	token2.Wait()
	return token2.Error()
}

func (b *Broker) Subscribe(ctx context.Context, channel string, handler brokers.Handler) (brokers.Subscription, error) {
	token := b.client.Subscribe(channel, 1, func(client mqtt.Client, mqttMsg mqtt.Message) {
		msg, err := models.ParseMessage(mqttMsg.Payload())
		if err != nil {
			// todo: be nice to have a way for user to deal with errors
			return
		}
		handler(msg)
	})
	token.Wait()
	if token.Error() != nil {
		return nil, token.Error()
	}
	return &Subscription{broker: b, channel: channel}, nil
}

func (b *Broker) Close() {
	// b.client.Close()
}

type Subscription struct {
	broker  *Broker
	channel string
}

func (s *Subscription) Unsubscribe() error {
	token := s.broker.client.Unsubscribe(s.channel)
	token.Wait()
	return token.Error()
}
