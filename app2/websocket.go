package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/treeder/async-toy-store/models"
	"github.com/treeder/async-toy-store/brokers"

)

var upgrader = websocket.Upgrader{} // use default options

// webSocketHandler will both listen to requests and pass them on to the correct channel
type webSocketHandler struct {
	c brokers.Broker
}

// I can see this being it's own "thing" that connects things together with one interface.
// It sits in front of and exposes the brokers.
// todo: what about auth?
func (h *webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// todo: this should check acceptable origins list
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade error:", err)
		return
	}
	log.Println("upgraded to websocket connection")
	defer c.Close()

	ctx := r.Context()
	sub, err := h.c.Subscribe(ctx, "orders_status", func(msg *models.Message) {
		fmt.Printf("wsHandler: Received a message on orders_status: %+v\n", msg)
		err2 := c.WriteJSON(msg)
		if err2 != nil {
			log.Print("WriteJSON error:", err)
		}
	})
	if err != nil {
		log.Print("error subscribing to orders_status:", err)
		return
	}
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = parseAndPublish(ctx, message)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		err = c.WriteMessage(mt, []byte(ack))
		if err != nil {
			log.Println("write error:", err)
			break
		}
	}
	sub.Unsubscribe()

}
