import { Component, OnInit } from '@angular/core';
import { BrokerService } from '../broker.service';
import { Order } from '../order';

@Component({
  selector: 'app-order-status',
  templateUrl: './order-status.component.html',
  styleUrls: ['./order-status.component.css']
})
export class OrderStatusComponent implements OnInit {

  order: Order;

  constructor(private broker: BrokerService) { }

  ngOnInit() {
    this.broker.subscribe("orders_status").subscribe(response => {
      if (response.channel != "orders_status") {  // todo: do this in the broker service
        return
      }
      console.log("orders_status message:", response);
      this.order = response.payload;
    })
  }

}
