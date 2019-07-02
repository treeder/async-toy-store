import { Component, OnInit } from '@angular/core';
import { FormGroup, FormControl } from '@angular/forms';
import { BrokerService } from '../broker.service';
import { Order } from '../order';

@Component({
  selector: 'app-checkout',
  templateUrl: './checkout.component.html',
  styleUrls: ['./checkout.component.css']
})
export class CheckoutComponent implements OnInit {

  checkoutForm = new FormGroup({
    firstName: new FormControl(''),
    lastName: new FormControl(''),
  });

  constructor(private broker: BrokerService) { }

  ngOnInit() {
  }

  onSubmit() {
    console.warn(this.checkoutForm.value);
    let order = new Order();
    order.id = "ABC";
    order.amount = 101.01;
    this.broker.publish("orders", order);
  }

}
