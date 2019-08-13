package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	nats "github.com/nats-io/nats.go"
	"github.com/rs/cors"
	"github.com/streadway/amqp"
	"github.com/treeder/async-toy-store/models"
)

const (
	ack = "{\"message\": \"ACK\"}"
)

var (
	natsClient *nats.EncodedConn
	mqttClient mqtt.Client
	amqpClient *amqp.Connection
)

func main() {

	// Connect to NATS
	// todo: move this into brokers/nats package, then use the /brokers interfaces here
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatalf("nats: %v", err)
	}
	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatalf("nats: %v", err)
	}
	defer c.Close()
	natsClient = c

	// Connect to MQTT
	// todo: move this into brokers/mqtt package, then use the /brokers interfaces here
	mqttClient = mqtt.NewClient(mqtt.NewClientOptions().AddBroker("ws://localhost:9005"))
	token := mqttClient.Connect()
	token.Wait()
	if token.Error() != nil {
		log.Fatalf("mqtt: %v", token.Error())
	}

	// Connect to RabbitMQ
	// More info: https://www.rabbitmq.com/tutorials/tutorial-one-go.html
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("rabbit: %v", err)
	}
	defer conn.Close()

	// Starts our REST/WebSockets API proxy so browsers can connect
	go startProxy(c)

	// Now subscribe
	var wg sync.WaitGroup
	wg.Add(1)
	c.Subscribe("orders", func(msg *models.Message) {
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

		// Put on next queue for others who are interested (ie: fulfillment)
		if err := c.Publish(paidChannel, msg2); err != nil {
			fmt.Printf("error publishing to %v: %v\n", paidChannel, err)
			return
		}
		msg2marshalled, err := json.Marshal(msg2)
		if err != nil {
			fmt.Printf("error marshalling msg2 for %v: %v\n", paidChannel, err)
			return
		}
		// also put on MQTT queue
		token2 := mqttClient.Publish(paidChannel, 1, false, msg2marshalled)
		token2.Wait()
		if token2.Error() != nil {
			log.Fatal(token2.Error())
		}
		// and why not RabbitMQ too? :)
		// todo: probably just want to get a channel once on init
		ch, err := conn.Channel()
		if err != nil {
			log.Fatal(err)
		}
		defer ch.Close()
		q, err := ch.QueueDeclare(
			paidChannel, // name
			false,       // durable
			false,       // delete when unused
			false,       // exclusive
			false,       // no-wait
			nil,         // arguments
		)
		if err != nil {
			log.Fatal(err)
		}

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        msg2marshalled,
			})

		// put a reply on reply channel
		if msg.ReplyChannel != "" {
			msg2 := &models.Message{Channel: msg.ReplyChannel, Payload: payload}
			if err := c.Publish(msg.ReplyChannel, msg2); err != nil {
				fmt.Println("error publishing to orders_status:", err)
				return
			}
		}
		fmt.Printf("order published to %v\n", paidChannel)
	})
	fmt.Println("Waiting for orders...")
	wg.Wait()

}

// This is a simple http proxy to nats, since there is almost no broker that seems to be able to work directly from the browser
// todo: move this into a separate utility that could do this for all brokers
func startProxy(c *nats.EncodedConn) {

	mux := http.NewServeMux()
	mux.Handle("/nats", &natsHandler{c: c})
	mux.Handle("/ws", &webSocketHandler{c: c})
	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

type natsHandler struct {
	c *nats.EncodedConn
}

func (h *natsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("nats called")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error:", err)
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	err = parseAndPublish(data)
	if err != nil {
		fmt.Println("error:", err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, ack)
}

func parseAndPublish(data []byte) error {
	msg, _, err := models.ParseMessage(data)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	if err := natsClient.Publish(msg.Channel, msg); err != nil {
		fmt.Println("error publishing:", err)
		return err
	}
	// Make sure the message goes through before we close
	natsClient.Flush()
	return nil
}
