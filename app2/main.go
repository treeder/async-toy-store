package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	nats "github.com/nats-io/nats.go"
	"github.com/rs/cors"
	"github.com/treeder/async-toy-store/models"
)

func main() {

	// todo: move this into brokers/nats package, then use the /brokers interfaces here
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	go startProxy(c)

	// Now subscribe
	var wg sync.WaitGroup
	wg.Add(1)
	c.Subscribe("orders", func(o *models.Order) {
		fmt.Printf("Received an order: %+v\n", o)

		// DO STUFF HERE

		// Enhance the order with payment information
		o.PaymentID = "paid123"

		// Put on next queue for others who are interested (ie: fulfillment)
		if err := c.Publish("orders_paid", o); err != nil {
			fmt.Println("error publishing to orders_paid:", err)
			return
		}
		fmt.Println("order published to orders_paid")
	})
	fmt.Println("Waiting for orders...")
	wg.Wait()

}

// This is a simple http proxy to nats, since there is almost no broker that seems to be able to work directly from the browser
// todo: move this into a separate utility that could do this for all brokers
func startProxy(c *nats.EncodedConn) {

	mux := http.NewServeMux()
	mux.HandleFunc("/nats", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("nats called")
		// b, _ := ioutil.ReadAll(r.Body)
		// fmt.Println("got:", string(b))
		order := &models.Order{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(order)
		if err != nil {
			// if err != io.EOF {
			fmt.Println("error:", err)
			http.Error(w, "Invalid format", http.StatusBadRequest)
			return
			// }
		}
		fmt.Printf("ORDER: %+v\n", order)
		if err := c.Publish("orders", &models.Order{ID: "123", Amount: 101.01}); err != nil {
			fmt.Println("error publishing:", err)
			http.Error(w, "Error publishing", http.StatusInternalServerError)
			return
		}
		// Make sure the message goes through before we close
		c.Flush()
		fmt.Fprintf(w, "{\"message\": \"ACK\"}")
	})
	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
