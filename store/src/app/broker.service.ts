import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { catchError, map, tap } from 'rxjs/operators';

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
  // client: MqttService;

  constructor(private http: HttpClient) {
    // this.connectToBroker();
  }

  async connectToBroker() {
    // try {
    //   this.client = await connect({ servers: ['nats://demo.nats.io:4222', 'tls://demo.nats.io:4443'] });
    //   // Do something with the connection
    // } catch (ex) {
    //   // handle the error
    // }

    // this.client = connect({ url: 'http://0.0.0.0:8081', json: false })

    // const cl = connect("ws://test.mosquitto.org:8080");
    // this.client = new AsyncClient(<IMqttClient>cl);

    // let s = new MqttService(MQTT_SERVICE_OPTIONS);
    // this.client = s;



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
    this.http.post<Order>("http://localhost:8080/nats", order)
    .pipe(
      // catchError(err => console.log("ERROR:", err))
    ).subscribe(response => {
      console.log("response:", response);
    })

  }

}
