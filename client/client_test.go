package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/marcell7/MQ/broker"
)

func TestPublisher(t *testing.T) {
	b := broker.New("127.0.0.1:3000")
	b.AddTopic("default")

	fmt.Println("Listening on port 3000")
	go b.Listen()

	time.Sleep(2 * time.Second)
	subscriber, err := NewSubscriber("127.0.0.1:3000")
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
	err = subscriber.Subscribe("default")
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	time.Sleep(2 * time.Second)
	publisher, err := NewPublisher("127.0.0.1:3000")
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
	err = publisher.Publish(`{"topic":"default","message":"Hello!"}`)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
	time.Sleep(5 * time.Second)
	for _, v := range b.Topics["default"].Subscriptions {
		if len(v.Queue) != 1 {
			t.Errorf("each subscription on default topic is expected to have 1 item in the queue got %d", len(v.Queue))
			return
		}
	}
}

func TestSubscriber(t *testing.T) {
	b := broker.New("127.0.0.1:3001")
	b.AddTopic("default")

	fmt.Println("Listening on port 3001")
	go b.Listen()
	time.Sleep(2 * time.Second)
	subscriber, err := NewSubscriber("127.0.0.1:3001")
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
	err = subscriber.Subscribe("default")
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
	time.Sleep(5 * time.Second)
	if len(b.Topics["default"].Subscriptions) != 1 {
		t.Errorf("Default topic should have 1 active subscription got %d", len(b.Topics["default"].Subscriptions))
	}
}

func TestSubscriberReceive(t *testing.T) {
	b := broker.New("127.0.0.1:3002")
	b.AddTopic("default")

	fmt.Println("Listening on port 3002")
	go b.Listen()

	time.Sleep(2 * time.Second)
	subscriber, err := NewSubscriber("127.0.0.1:3002")
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
	publisher, err := NewPublisher("127.0.0.1:3002")
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
	err = publisher.Publish(`{"topic":"default","message":"Hello World!"}`)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
	time.Sleep(3 * time.Second)
	// -------------------------

	// Receive current message from the topic
	msg, _ := subscriber.Receive("default")
	if msg.Payload.Message != "Hello World!" {
		t.Errorf("Expected to receive Hello World! got %s", msg.Payload.Message)
	}

	for _, subscription := range b.Topics["default"].Subscriptions {
		if len(subscription.Queue) != 0 {
			t.Errorf("Expected 0 items in a queue got %d", len(subscription.Queue))
		}
	}

}

func TestReceiveFromEmptyQueue(t *testing.T) {
	b := broker.New("127.0.0.1:3002")
	b.AddTopic("default")

	fmt.Println("Listening on port 3002")
	go b.Listen()

	time.Sleep(2 * time.Second)
	subscriber, err := NewSubscriber("127.0.0.1:3002")
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
	err = subscriber.Subscribe("default")
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
	time.Sleep(3 * time.Second)

	// Receive current message from the topic
	_, err = subscriber.Receive("default")
	if err != nil {
		return
	} else {
		t.Errorf("Expected error got none")
	}
}
