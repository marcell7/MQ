package broker

import (
	"errors"
	"sync"
)

type Subscription struct {
	id         string       // id
	subscriber *Subscriber  // subscriber
	mu         sync.RWMutex // mutex for reading and writing to the queue
	Queue      []*Item      // queue that holds published items
}

// Constructor for Subscription struct
func newSubscription(id string, subscriber *Subscriber) *Subscription {
	return &Subscription{
		id:         id,
		subscriber: subscriber,
	}
}

// Add item to the queue
func (s *Subscription) addToQueue(item *Item) {
	s.mu.Lock()
	s.Queue = append(s.Queue, item)
	s.mu.Unlock()
}

// Take the item out of the queue and return it
func (s *Subscription) popOut() (*Item, error) {
	s.mu.Lock()
	if len(s.Queue) == 0 {
		return nil, errors.New("queue is empty")
	}
	currentItem := s.Queue[len(s.Queue)-1]
	s.Queue = s.Queue[:len(s.Queue)-1]
	s.mu.Unlock()
	return currentItem, nil
}
