package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/treeder/async-toy-store/models"
)

func main() {

	// todo: move this into brokers/mqtt package, then use the /brokers interfaces here
	c := mqtt.NewClient(mqtt.NewClientOptions().AddBroker("ws://localhost:9005"))
	token := c.Connect()
	token.Wait()
	if token.Error() != nil {
		log.Fatal(token.Error())
	}

	var wg sync.WaitGroup
	wg.Add(1)
	// todo: this is just grabbing it from the browser order, should grab it from payment processing
	c.Subscribe("orders", 1, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("app3: Received an order: %+v\n", msg)
		order := &models.Order{}
		err := json.Unmarshal(msg.Payload(), order)
		if err != nil {
			log.Println("error unmarshalling order", err)
		}
		fmt.Printf("app3: order: %+v", order)

		// DO STUFF HERE

		// TODO: enhance the order with payment info or tracking info, order.PaymentID or order.TrackingNumber
		// TODO: then put on "orders2" queue

		// wg.Done()
	})
	fmt.Println("app3: Waiting for orders...")
	wg.Wait()

}
