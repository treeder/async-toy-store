# async-toy-store

Demo for AsyncAPI

Parts:

* a UI at /store to checkout and send the first "order" event/message
* goapp represents payment processing. Subscribes to order events from browser and publishes enhanced event with Payment ID.
* jsapp is fulfillment processing. Subscribes to order events from goapp and publishes enhanced event with tracking number.

The client (browser UI) needs to know what the payment processor accepts, so it uses the payment processors AsyncAPI spec,
much like it would use it's REST spec. From this it knows what server to send order events to and in what format. In this
demo that is the `orders` channel. It does not need its own AsyncAPI spec.

The payment processor subscribes to order events, processes payments, then publishes the event with extra payment info to another channel,
`orders_paid`.

Fulfillment is interested in the events on the `orders_paid` channel so it subscribes to them. For this it reads the AsyncAPI spec from
the payment processor at path `{URL_TO_SPEC}/orders_paid/publish` (JS Link?) since that is who defines both whether the events will go and the format they will be in. Fulfillment then publishes the order event with an added tracking number to `orders_shipped` channel.

## Walkthrough

### Generating code

Note: this part doesn't work yet, need more code generators.

TODO: AsyncAPI needs to support more languages for generating code.
TODO: AsyncAPI should also support various message brokers while generating code, if it doesn't already.

```sh
# should be once javascript is supported in new one: docker run --rm -it -v $PWD:/app -w /app treeder/asyncapi-gen node cli -o output orders.yaml javascript
# but for now:
docker run --rm -it -v $PWD:/app -w /app treeder/asyncapi-gen node ac LOOK IT UP
```

### Start Services

Start nats server:

```sh
docker run --rm -it -d --name nats -p 4222:4222 -p 6222:6222 -p 8222:8222 nats
```

Start payment processor app:

```sh
cd app1
make run
```

Start UI:

```sh
cd store
make run
```

Open UI at http://localhost:4200

## For Figuring out Later

### Nats

Unfortunately most message brockers don't work nicely directly from a browser, some requiring a websocket to tcp proxy and most of the client libraries expect Node libraries
so they don't work in the browser. I added a simple REST proxy that the UI talks to for now.
