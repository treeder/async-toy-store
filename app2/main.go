package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/rs/cors"
	"github.com/treeder/async-toy-store/brokers"
	"github.com/treeder/async-toy-store/brokers/auto"
	"github.com/treeder/async-toy-store/models"
)

const (
	ack = "{\"message\": \"ACK\"}"
)

var (
	natsClient brokers.Broker
	mqttClient brokers.Broker
	amqpClient brokers.Broker
)

func main() {

	ctx := context.Background()
	var err error

	// Connect to NATS
	natsClient, err = auto.Connect(ctx, "nats://localhost:4222")
	if err != nil {
		log.Fatalf("nats: %v", err)
	}
	defer natsClient.Close()

	// Connect to MQTT
	mqttClient, err = auto.Connect(ctx, "mqtt+ws://localhost:9005")
	if err != nil {
		log.Fatalf("mqtt: %v", err)
	}
	defer mqttClient.Close()

	// Connect to RabbitMQ
	amqpClient, err = auto.Connect(ctx, "amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("amqp: %v", err)
	}
	defer amqpClient.Close()

	// Starts our REST/WebSockets API proxy so browsers can connect
	go startProxy(natsClient)

	// Now subscribe
	var wg sync.WaitGroup
	wg.Add(1)
	natsClient.Subscribe(ctx, "orders", func(msg *models.Message) {
		fmt.Printf("Received a message on orders channel: %+v\n", msg)
		o, err := models.ParseOrder(msg.Payload)
		if err != nil {
			fmt.Println("error parsing order from msg:", err)
			return
		}

		fmt.Printf("Received an order: %+v\n", o)

		// DO STUFF HERE
		time.Sleep(1 * time.Second)

		// Enhance the order with payment information
		o.PaymentID = "paid123"
		o.Status = "Payment successful"

		paidChannel := "orders_paid"

		payload, err := json.Marshal(o)
		if err != nil {
			fmt.Printf("error marshalling order for %v: %v\n", paidChannel, err)
			return
		}
		msg2 := &models.Message{Channel: paidChannel, ReplyChannel: msg.ReplyChannel, Payload: payload}
		if err != nil {
			fmt.Printf("error marshalling msg2 for %v: %v\n", paidChannel, err)
			return
		}

		// Put on next queue for others who are interested (ie: fulfillment)
		if err := natsClient.Publish(ctx, paidChannel, msg2); err != nil {
			fmt.Printf("error publishing to %v: %v\n", paidChannel, err)
			return
		}

		// also put on MQTT queue
		if err := mqttClient.Publish(ctx, paidChannel, msg2); err != nil {
			fmt.Printf("error publishing to %v: %v\n", paidChannel, err)
			return
		}

		// and why not RabbitMQ too? :)
		if err := amqpClient.Publish(ctx, paidChannel, msg2); err != nil {
			fmt.Printf("error publishing to %v: %v\n", paidChannel, err)
			return
		}

		// put a reply on reply channel
		if msg.ReplyChannel != "" {
			msg2 := &models.Message{Channel: msg.ReplyChannel, Payload: payload}
			if err := natsClient.Publish(ctx, msg.ReplyChannel, msg2); err != nil {
				fmt.Println("error publishing to orders_status:", err)
				return
			}
		}
		fmt.Printf("order published to %v\n", paidChannel)
	})
	fmt.Println("app2: Waiting for orders...")
	wg.Wait()

}

// This is a simple http proxy to nats, since there is almost no broker that seems to be able to work directly from the browser
// todo: move this into a separate utility that could do this for all brokers
func startProxy(c brokers.Broker) {

	mux := http.NewServeMux()
	// mux.Handle("/nats", &natsHandler{c: c})
	mux.Handle("/ws", &webSocketHandler{c: c})
	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// this is a REST proxy endpoint
type natsHandler struct {
	c brokers.Broker
}

func (h *natsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Println("nats called")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error:", err)
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	err = parseAndPublish(ctx, data)
	if err != nil {
		fmt.Println("error:", err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, ack)
}

func parseAndPublish(ctx context.Context, data []byte) error {
	msg, err := models.ParseMessage(data)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	if err := natsClient.Publish(ctx, msg.Channel, msg); err != nil {
		fmt.Println("error publishing:", err)
		return err
	}
	// Make sure the message goes through before we close
	// natsClient.Flush()w
	return nil
}
