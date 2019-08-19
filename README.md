# async-toy-store

Demo for AsyncAPI

Parts:

* app1: a UI at /store to checkout and send the first "order" event/message to a channel that app2 is subscribed to.
* app2: is payment processing. Subscribes to order events from browser and publishes enhanced event with Payment ID.
* app3: is fulfillment. Subscribes to order events produced from app2 and publishes enhanced event with tracking number.
* muleapp: is processing events from the orders_paid channel, enriching them with customer data and putting them in the orders_status channel.

The client (browser UI) needs to know what the payment processor accepts, so it uses the payment processors AsyncAPI spec,
much like it would use it's REST spec. From this it knows what server to send order events to and in what format. In this
demo that is the `orders` channel. It does not need its own AsyncAPI spec.

The payment processor subscribes to order events, processes payments, then publishes the event with extra payment info to another channel,
`orders_paid`.

Fulfillment is interested in the events on the `orders_paid` channel so it subscribes to them. For this it reads the AsyncAPI spec from
the payment processor at path `{URL_TO_SPEC}/orders_paid/publish` (JS Link?) since that is who defines both whether the events will go and the format they will be in. Fulfillment then publishes the order event with an added tracking number to `orders_shipped` channel.

## Walkthrough

### Generating code

* TODO: AsyncAPI needs to support more languages for generating code.
* TODO: AsyncAPI should also support various message brokers while generating code.

```sh
# should be once javascript is supported in new one: docker run --rm -it -v $PWD:/app -w /app treeder/asyncapi-gen node cli -o output orders.yaml javascript
# but for now:
docker run --rm -it -v $PWD:/app -w /app treeder/asyncapi-gen node ac LOOK IT UP
```

### Start Services

Start Nats server:

```sh
docker run --rm -it -d --name nats -p 4222:4222 -p 6222:6222 -p 8222:8222 nats
```

Start MQTT (Mosquitto) server:

```sh
docker run --rm -it -d --name mosquitto -p 1883:1883 -p 9005:9005 -v $PWD/mosquitto.conf:/mosquitto/config/mosquitto.conf eclipse-mosquitto
```

Start RabbitMQ server:

```sh
docker run --rm -it -d --name rabbit -p 5672:5672 --hostname my-rabbit rabbitmq:3-alpine
```

Start all apps in different consoles:

```sh
cd app1
make run
```

```sh
cd app2
make run
```

```sh
cd app3
make run
```


```sh
cd muleApp
sh buildAndRun.sh```

Open UI at http://localhost:4200

## For Figuring out Later

Unfortunately most message brokers don't work nicely directly from a browser, some requiring a websocket to tcp proxy and most of the client libraries expect Node libraries so they don't work in the browser. I added a simple REST proxy that the UI can talk to.

## TODO

* [ ] Maybe make the Message object we're passing around a CloudEvent object instead?
