import { Injectable } from '@angular/core';
import { connect, Client, NatsConnectionOptions, Payload } from 'ts-nats';

@Injectable({
  providedIn: 'root'
})
export class NatsService {

  nc: Client;

  constructor() {
    this.connectToNats();
  }

  async connectToNats()  {
    try {
      this.nc = await connect({ servers: ['nats://demo.nats.io:4222', 'tls://demo.nats.io:4443'] });
      // Do something with the connection
    } catch (ex) {
      // handle the error
    }
  }

  publish(channel: string, data: any) {
    this.nc.publish(channel, data)
  }
  
}
