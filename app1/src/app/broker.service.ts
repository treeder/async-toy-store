import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {Client as MqttClient, Message as MqttMessage} from 'paho-mqtt'
import { webSocket, WebSocketSubject } from "rxjs/webSocket";

import { Order } from './order';
import { Message } from './message';
import { Subscription, Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class BrokerService {

  client: MqttClient;
  wsSubject: WebSocketSubject<Message<Order>>;

  constructor(private http: HttpClient) {
    this.connectToBroker();
  }

  connectToBroker() {
    console.log("connecting to MQTT")
    this.client = new MqttClient(location.hostname, Number(9005), "clientId-app1");
    // set callback handlers
    this.client.onConnectionLost = this.onConnectionLost;
    this.client.onMessageArrived = this.onMessageArrived;

    // connect the client
    this.client.connect({ onSuccess: this.onConnect.bind(this) });

    console.log("also connecting to websocket interface");
    this.wsSubject = webSocket('ws://localhost:8080/ws');
    this.wsSubject.subscribe(response => {
      console.log("ws response:", response);
    })
  }

  // called when the client connects 
  onConnect() {
    // Once a connection has been made, make a subscription and send a message.
    console.log("onConnect");
    this.client.subscribe("World");
    let message = new MqttMessage("Hello");
    message.destinationName = "World";
    this.client.send(message);
  }

  // called when the client loses its connection
  onConnectionLost(responseObject) {
    if (responseObject.errorCode !== 0) {
      console.log("onConnectionLost:" + responseObject.errorMessage);
    }
  }

  // called when a message arrives
  onMessageArrived(message) {
    console.log("onMessageArrived:" + message.payloadString);
  }

  publish(channel: string, order: Order) {

    let msg = new Message<Order>();
    msg.channel = channel;
    msg.reply_channel = "orders_status";
    msg.payload = order;
    
    // This one uses Nats proxy
    order.comment = "using nats"
    this.http.post<Order>("http://localhost:8080/nats", msg)
      .pipe(
        // catchError(err => console.log("ERROR:", err))
      ).subscribe(response => {
        console.log("response:", response);
      })

    // This one uses Paho MQTT, seems to be the only one that works in a browser
    order.comment = "using mqtt"
    let message = new MqttMessage(JSON.stringify(msg));
    message.destinationName = "orders";
    console.log("and mqtt orders")
    this.client.send(message);

    // And finally one using our websockets interface
    // todo: should wrap the object in a Message object that states the channel
    order.comment = "websocket"
    this.wsSubject.next(msg);

  }


  subscribe(channel: string): Observable<Message<Order>> {
    // todo: try multiplexing rxjs filter here to filter by channel: https://rxjs-dev.firebaseapp.com/api/webSocket/webSocket
    return this.wsSubject;
  }

}
