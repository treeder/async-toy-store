package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/treeder/async-toy-store/brokers"
	"github.com/treeder/async-toy-store/brokers/auto"
	"github.com/treeder/async-toy-store/models"
)

var (
	natsClient brokers.Broker
	mqttClient brokers.Broker
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

	var wg sync.WaitGroup
	wg.Add(1)
	// todo: this is just grabbing it from the browser order, should grab it from payment processing
	mqttClient.Subscribe(ctx, "orders_paid", func(msg *models.Message) {
		fmt.Printf("app3: Received an order: %+v\n", msg)
		// msg, err := models.ParseMessage(msg.Payload())
		// if err != nil {
		// 	fmt.Println("error:", err)
		// 	return
		// }
		order, err := models.ParseOrder(msg.Payload)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Printf("app3: order: %+v", order)

		// DO STUFF HERE
		time.Sleep(2 * time.Second)
		order.TrackingID = "123456"
		order.Status = "Shipped"

		// enhance the order with payment info or tracking info, order.PaymentID or order.TrackingNumber
		if msg.ReplyChannel != "" {
			payload, err := json.Marshal(order)
			if err != nil {
				fmt.Println("error marshalling order:", err)
				return
			}
			msg2 := &models.Message{Channel: msg.ReplyChannel, Payload: payload}
			// this should end up back at the browser
			// TODO: there's not enough information in the message to know which broker to send this to...
			if err := natsClient.Publish(ctx, "orders_status", msg2); err != nil {
				fmt.Println("error publishing to orders_status:", err)
			}
		}
		// wg.Done()
	})
	fmt.Println("app3: Waiting for orders_paid...")
	wg.Wait()

}
