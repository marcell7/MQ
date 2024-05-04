package broker

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/marcell7/teenypubsub/protocol"
)

type Server interface {
	Listen() error // Starts listening and serving
	Stop() error   // Stops the server
}

// Implements the Server interface
type Broker struct {
	listenAddr  string                   // Address the broker is listening on
	listener    net.Listener             // Listener object
	mu          sync.RWMutex             // Mutex for adding and removing to and from publishers, subscribers and Topics maps
	publishers  map[string]*Publisher    // Map that stores publishers registered on the broker - {"<publisher_id":"<Publisher>"}
	subscribers map[string]*Subscriber   // Map that stores subscribers registered on the broker - {"<subscriber_id":"<Subscriber>"}
	Topics      map[string]*DefaultTopic // Map that stores topics on the broker - {"<topic_id>":"<DefaultTopic>"}
	protocol    protocol.Protocol        // Protocol object holding the methods required for decoding/encoding data sent to and from the broker

	exitCh chan struct{} // Channel used for signaling when to exit the server. Used for manually stopping the server
}

// Constructor for the Broker struct
func New(listenAddr string) *Broker {
	return &Broker{
		listenAddr:  listenAddr,
		publishers:  make(map[string]*Publisher),
		subscribers: make(map[string]*Subscriber),
		protocol:    new(protocol.DefaultProtocol),
		Topics:      make(map[string]*DefaultTopic),
		exitCh:      make(chan struct{}),
	}
}

func (b *Broker) Listen() error {
	listener, err := net.Listen("tcp", b.listenAddr)
	if err != nil {
		return err
	}
	b.listener = listener
	go b.startAcceptingConnections()

	return nil
}

func (b *Broker) Stop() {
	log.Println("Stopping the broker...")
	close(b.exitCh)
	b.listener.Close()

}

func (b *Broker) AddTopic(name string) {
	id := generateId()
	topic := newDefaultTopic(id, name)
	b.Topics[name] = topic
}

func (b *Broker) startAcceptingConnections() error {
	for {
		conn, err := b.listener.Accept()
		if err != nil {
			select {
			case <-b.exitCh:
				// Server was manually stopped
				return nil
			default:
				fmt.Printf("Error %s", err)
				continue
			}
		}
		go b.handleConnection(conn)
	}
}

func (b *Broker) handleConnection(conn net.Conn) error {
	var publisher *Publisher
	var subscriber *Subscriber
	defer func() {
		if publisher != nil {
			b.removeClient(publisher)
		}
		if subscriber != nil {
			b.removeClient(subscriber)
		}
		fmt.Println("Dropping a connection with client")
		conn.Close()
	}()
	clientId := generateId()
	reader := bufio.NewReader(conn)
	for {
		data, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				// No data in the reader
				return err
			} else {
				fmt.Printf("error: %s\n", err)
				return err
			}

		}
		msg := &protocol.DefaultMessage{}
		if err := b.protocol.Decode(msg, data); err != nil {
			return err
		}

		switch msg.Command {
		case protocol.CMD_PUBREG:
			publisher = newPublisher(clientId, conn)
			b.addClient(publisher)
			publisher.sendOk()
		case protocol.CMD_SUBREG:
			subscriber = newSubscriber(clientId, conn)
			b.addClient(subscriber)
			subscriber.sendOk()
		case protocol.CMD_PUB:
			if _, ok := b.publishers[clientId]; ok {
				item := newItem(msg.Id, msg.Payload.Message)
				if err := b.Topics[msg.Payload.Topic].addItem(item); err != nil {
					publisher.sendError("no active subscriptions")
				}
				err = publisher.sendOk()
				if err != nil {
					return err
				}
			} else {
				publisher.sendError("must be registered as a publisher")
				return errors.New("must be registered as a publisher")
			}
		case protocol.CMD_SUB:
			if subscriber, ok := b.subscribers[clientId]; ok {
				b.Topics[msg.Payload.Topic].addSubscription(clientId, subscriber)
				err = subscriber.sendOk()
				if err != nil {
					return err
				}
			} else {
				subscriber.sendError("must be registered as a subscriber")
				return errors.New("must be registered as a subscriber")
			}
		case protocol.CMD_RECV:
			if subscription, ok := b.Topics[msg.Payload.Topic].Subscriptions[clientId]; ok {
				currentItem, err := subscription.popOut()
				if err != nil {
					subscriber.sendError("no items in the queue")
					return err
				}
				subscriber.sendResp(currentItem)
			}
		}
	}
}

func (b *Broker) addClient(client Client) {
	b.mu.Lock()
	switch c := client.(type) {
	case *Publisher:
		b.publishers[c.id] = c
	case *Subscriber:
		b.subscribers[c.id] = c
	}
	b.mu.Unlock()
}

func (b *Broker) removeClient(client Client) {
	b.mu.Lock()
	switch c := client.(type) {
	case *Publisher:
		delete(b.publishers, c.id)
	case *Subscriber:
		delete(b.subscribers, c.id)
	}
	b.mu.Unlock()
}
