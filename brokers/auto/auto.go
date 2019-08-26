package auto

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/treeder/async-toy-store/brokers"
	"github.com/treeder/async-toy-store/brokers/mqtt"
	"github.com/treeder/async-toy-store/brokers/nats"
	"github.com/treeder/async-toy-store/brokers/rabbit"
)

// Given a URL, connects to the correct broker
func Connect(ctx context.Context, urlStr string) (brokers.Broker, error) {

	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	s := url.Scheme
	s = strings.Split(s, "+")[0]
	switch s {
	case "nats":
		return nats.Connect(ctx, urlStr)
	case "mqtt":
		return mqtt.Connect(ctx, urlStr)
	case "amqp":
		return rabbit.Connect(ctx, urlStr)
	}
	return nil, fmt.Errorf("No broker found for %v", url.Scheme)

}
