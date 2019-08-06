package brokers

type Broker interface {
	Publish(topic string, message Message) error
	Subscribe(topic string, handler Handler)
}

type Message interface {
	Payload() []byte
}

type Handler func(Message)
