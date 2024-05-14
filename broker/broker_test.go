package broker

import (
	"fmt"
	"testing"
	"time"

	"github.com/marcell7/MQ/client"
)

func TestDeleteSubscriptionOnSubscriberDisconnect(t *testing.T) {
	b := New("127.0.0.1:3002")
	b.AddTopic("default")

	fmt.Println("Listening on port 3002")
	go b.Listen()

	time.Sleep(2 * time.Second)
	subscriber, err := client.NewSubscriber("127.0.0.1:3002")
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
	err = subscriber.Subscribe("default")
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	// Init publisher
	time.Sleep(2 * time.Second)
	publisher, err := client.NewPublisher("127.0.0.1:3002")
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
	err = publisher.Publish(`{"topic":"default","message":"Hello World!"}`)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
	time.Sleep(2 * time.Second)

	subscriber.Close()

	if len(b.Topics["default"].Subscriptions) != 0 {
		t.Errorf("Expected 0 subscriptions got %d", len(b.Topics["default"].Subscriptions))
	}
}
