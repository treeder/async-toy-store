package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	nats "github.com/nats-io/nats.go"
	"github.com/treeder/async-toy-store/models"
)

var (
	natsClient *nats.EncodedConn
	mqttClient mqtt.Client
)

func main() {

	// Connect to NATS
	// todo: move this into brokers/nats package, then use the /brokers interfaces here
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	natsClient, err = nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	defer natsClient.Close()

	// Connect to MQTT
	// todo: move this into brokers/mqtt package, then use the /brokers interfaces here
	mqttClient = mqtt.NewClient(mqtt.NewClientOptions().AddBroker("ws://localhost:9005"))
	token := mqttClient.Connect()
	token.Wait()
	if token.Error() != nil {
		log.Fatal(token.Error())
	}

	var wg sync.WaitGroup
	wg.Add(1)
	// todo: this is just grabbing it from the browser order, should grab it from payment processing
	mqttClient.Subscribe("orders_paid", 1, func(client mqtt.Client, mqttMsg mqtt.Message) {
		fmt.Printf("app3: Received an order: %+v\n", mqttMsg)
		msg, order, err := models.ParseMessage(mqttMsg.Payload())
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
			if err := natsClient.Publish("orders_status", msg2); err != nil {
				fmt.Println("error publishing to orders_status:", err)
			}
		}

		// wg.Done()
	})
	fmt.Println("app3: Waiting for orders_paid...")
	wg.Wait()

}
