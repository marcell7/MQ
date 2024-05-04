package broker

import (
	"errors"
	"sync"
)

type Topic interface {
	addItem(*Item) error                 // Adds an item to the topics's queue
	addSubscription(string, *Subscriber) // Adds a subscription to the topic
}

// Implements Topic interface
// Broker can have multiple topics. Each topic can have multiple subscriptions. Each subscription has a subscriber and queue that subscriber fetches items/messages from.
type DefaultTopic struct {
	id            string                   // id of the topic
	name          string                   // name of the topic
	mu            sync.RWMutex             // mutex for modifying the subscription map
	Subscriptions map[string]*Subscription // Map storing all subscriptions for that topic. Map key is the subscriber's id
}

// Constructor for the DefaultTopic struct
func newDefaultTopic(id string, name string) *DefaultTopic {
	return &DefaultTopic{
		id:            id,
		name:          name,
		Subscriptions: make(map[string]*Subscription),
	}
}

func (dt *DefaultTopic) addItem(item *Item) error {
	dt.mu.RLock()
	if len(dt.Subscriptions) == 0 {
		return errors.New("no active subscriptions on this topic")
	}
	dt.mu.RUnlock()
	// Each topic can have multiple subscriptions - one for each subscriber of that topic.
	// Add item to every queue in these subscriptions
	for _, subscription := range dt.Subscriptions {
		subscription.addToQueue(item)
	}

	return nil
}

func (dt *DefaultTopic) addSubscription(id string, subscriber *Subscriber) {
	dt.mu.Lock()
	// Change to &Subscription{id:id, subscriber:subscriber}
	subscription := newSubscription(id, subscriber)
	dt.Subscriptions[id] = subscription
	dt.mu.Unlock()
}

// Item struct that represents the data that is stored in the topic's queue
type Item struct {
	Id   string // id of the publish message
	Data string // User-provided data
}

// Constructor for the Item struct
func newItem(id string, data string) *Item {
	return &Item{
		Id:   id,
		Data: data,
	}
}
