# MQ

**!!! This is just an experimental message queue and should not be used in any serious projects !!!**

Meant for experimentation only.
___
```txt
                BROKER
              | topic1 | -> subscribers
publisher ->  |        |
              | topic2 | -> subscribers
```

## Usage

Run a broker:

```go
// Create a broker
b := broker.New("127.0.0.1:3000")
// Create a topic named default
b.AddTopic("default")
// Start listening and serving
go b.Listen()
```

Connect a subscriber that subscribes to the default topic

```go
// Connect a subscriber
subscriber, err := client.NewSubscriber("127.0.0.1:3000")
// Subscribe to the default topic
err := subscriber.Subscribe("default")
// Grab the message from the default topic
msg, err := subscriber.Receive("default")
```

Connect a publisher that publishes a message to the default topic

```go
// Connect a publisher
publisher, err := client.NewPublisher("127.0.0.1:3000")
// Publish a message to the default topic
err = publisher.Publish(`{"topic":"default","message":"Hello World!"}`)
```
