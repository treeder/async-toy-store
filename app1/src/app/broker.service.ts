import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {Client as MqttClient, Message as MqttMessage} from 'paho-mqtt'

import { Order } from './order';

// NONE OF THIS STUFF WORKS IN THE BROWSER... argggh
// import { connect, Client } from 'ws-nats';
// const MQTT = require("async-mqtt");
// import { connect, MqttClient } from 'mqtt';
// import { AsyncClient, IMqttClient } from 'async-mqtt';
// import {
//   IMqttMessage,
//   MqttModule,
//   IMqttServiceOptions,
//   MqttService
// } from 'ngx-mqtt';


// const MQTT_SERVICE_OPTIONS: IMqttServiceOptions = {
//   hostname: 'test.mosquitto.org',
//   port: 8080,
//   // path: '/mqtt'
// };

@Injectable({
  providedIn: 'root'
})
export class BrokerService {

  // client: Client;
  // client: AsyncClient;
  client: MqttClient;

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

  publish(channel: string, data: any) {
    // try {
    //   this.client.publish(channel, data);
    // } catch (err) {
    //   console.log("ERROR", err)
    // }
    // asyncClient.publish("foo/bar", "baz").then(() => {
    // 	console.log("We async now");
    // 	return asyncClient.end();
    // });

    let order = new Order();
    order.id = "ABC";
    order.amount = 101.01;
    // This one uses Nats proxy
    order.comment = "using nats"
    this.http.post<Order>("http://localhost:8080/nats", order)
      .pipe(
        // catchError(err => console.log("ERROR:", err))
      ).subscribe(response => {
        console.log("response:", response);
      })

    // This one uses Paho MQTT, seems to be the only one that works in a browser
    order.comment = "using mqtt"
    let message = new MqttMessage(JSON.stringify(order));
    message.destinationName = "orders";
    console.log("and mqtt orders")
    this.client.send(message);

  }

}
