# async-toy-store

Demo for AsyncAPI

Pieces

* a UI at /store to checkout and send the first "order" event/message
* goapp receives event (these can be built out to do specific things like one payment processing)
* jsapp receives event

## Walkthrough

Generate code:

```sh
# docker run --rm -it -p 8080:8080 -v $PWD:/app -w /app treeder/asyncapi-gen node cli -o output test/docs/streetlights.yml html
docker run --rm -it -v $PWD:/app -w /app treeder/asyncapi-gen node cli -o output orders.yaml javascript
```

