package rabbit

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
	"github.com/treeder/async-toy-store/brokers"
	"github.com/treeder/async-toy-store/models"
)

func Connect(ctx context.Context, urlStr string) (*Broker, error) {
	// More info: https://www.rabbitmq.com/tutorials/tutorial-one-go.html
	conn, err := amqp.Dial(urlStr)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	b := &Broker{conn: conn, ch: ch}
	return b, nil
}

type Broker struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func (b *Broker) Publish(ctx context.Context, channel string, msg *models.Message) error {
	// todo: should we cache the queues so we don't do this on every publish?
	q, err := b.ch.QueueDeclare(
		channel, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return err
	}
	msgMarshalled, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = b.ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msgMarshalled,
		})
	return err
}
func (b *Broker) Subscribe(ctx context.Context, channel string, handler brokers.Handler) (brokers.Subscription, error) {
	q, err := b.ch.QueueDeclare(
		channel, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, err
	}

	s := &Subscription{q: q, handler: handler}
	go s.start()
	return s, nil
}
func (b *Broker) Close() {
	b.ch.Close()
	b.conn.Close()

}

type Subscription struct {
	ch      *amqp.Channel
	q       amqp.Queue
	handler brokers.Handler
	msgs    chan *amqp.Delivery
}

func (s *Subscription) start() {
	msgs, err := s.ch.Consume(
		s.q.Name, // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		fmt.Println("ERROR starting AMQP Subscription:", err)
		return
	}

	for d := range msgs {
		// log.Printf("Received a message: %s", d.Body)
		msg, err := models.ParseMessage(d.Body)
		if err != nil {
			fmt.Printf("SUBSCRIBE ERROR: %v\n", err)
			// todo: be nice to have a way for user to deal with errors
			continue
		}
		s.handler(msg)
	}
}

func (s *Subscription) Unsubscribe() error {
	close(s.msgs)
	return nil
}
