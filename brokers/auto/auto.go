package auto

import(
	"context"
	"net/url"
	"fmt"

	"github.com/treeder/async-toy-store/brokers"
	"github.com/treeder/async-toy-store/brokers/nats"
)

// Given a URL, connects to the correct broker
func Connect(ctx context.Context, urlStr string) (brokers.Broker, error) {

	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	switch url.Scheme{
	case "nats":
return nats.Connect(ctx, urlStr)

	}
	return nil, fmt.Errorf("No broker found for %v", url.Scheme)


}