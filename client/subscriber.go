package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/marcell7/MQ/protocol"
)

type Subscriber interface {
	Subscribe(string) error                           // Subscribes to the user provided topic
	Receive(string) (*protocol.DefaultMessage, error) // Receive the last message from the topic queue (FIFO style)
	Close() error                                     // Closes the connection
	start() error                                     // Starts listening for incoming messages
	connect() error                                   // Connects the client to the broker (tcp server)
	register() error                                  // Registers the client as a subscriber on the broker
}

// Implements Subscriber interface
type DefaultSubscriber struct {
	addr       string                        // Address of the broker
	conn       net.Conn                      // Connection of that specific subscriber - allows for writing and receiving messages from / to the broker
	protocol   protocol.Protocol             // Protocol instance for decoding messages
	okCh       chan struct{}                 // Channel for signaling succesfully processed messages
	errCh      chan *protocol.DefaultMessage // Channel for errors encountered on the broker
	receiverCh chan *protocol.DefaultMessage // Channel used for receiving messages from the broker
}

// Constructor for the DefaultSubscriber struct
func NewSubscriber(addr string) (*DefaultSubscriber, error) {
	ds := &DefaultSubscriber{
		addr:       addr,
		protocol:   new(protocol.DefaultProtocol),
		okCh:       make(chan struct{}),
		errCh:      make(chan *protocol.DefaultMessage),
		receiverCh: make(chan *protocol.DefaultMessage),
	}
	if err := ds.connect(); err != nil {
		return nil, err
	}

	go ds.start()

	if err := ds.register(); err != nil {
		return nil, err
	}
	return ds, nil
}

func (ds *DefaultSubscriber) Subscribe(topic string) error {
	payload := fmt.Sprintf(`{"topic":"%s"}`, topic)
	if _, err := ds.conn.Write([]byte(fmt.Sprintf("SUB %s\n", payload))); err != nil {
		return err
	}
	select {
	case <-ds.okCh:
		// If subscribe message was processed succesfully
		return nil
	case errMsg := <-ds.errCh:
		// If subscribe message was NOT processed succesfully (received error from the broker)
		return errors.New(errMsg.Payload.Error)
	}
}

func (ds *DefaultSubscriber) Receive(topic string) (*protocol.DefaultMessage, error) {
	payload := fmt.Sprintf(`{"topic":"%s"}`, topic)
	if _, err := ds.conn.Write([]byte(fmt.Sprintf("RECV %s\n", payload))); err != nil {
		return nil, err
	}
	select {
	case msg := <-ds.receiverCh:
		// Received item/message from the topic queue
		return msg, nil
	case errMsg := <-ds.errCh:
		// Error received from the broker
		return nil, errors.New(errMsg.Payload.Error)
	}
}

func (ds *DefaultSubscriber) Close() error {
	if err := ds.conn.Close(); err != nil {
		return err
	}
	return nil
}

func (ds *DefaultSubscriber) start() error {
	reader := bufio.NewReader(ds.conn)
	for {
		data, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return err
			} else {
				fmt.Printf("error: %s\n", err)
				return err
			}
		}
		msg := &protocol.DefaultMessage{}
		if err := ds.protocol.Decode(msg, data); err != nil {
			fmt.Println("Error decoding")
			return err
		}
		switch msg.Command {
		case protocol.CMD_OK:
			ds.okCh <- struct{}{}
		case protocol.CMD_ERROR:
			ds.errCh <- msg
		case protocol.CMD_RESP:
			ds.receiverCh <- msg
		}

	}
}

func (ds *DefaultSubscriber) connect() error {
	conn, err := net.Dial("tcp", ds.addr)
	if err != nil {
		return err
	}
	ds.conn = conn
	return nil
}

func (ds *DefaultSubscriber) register() error {
	if _, err := ds.conn.Write([]byte("SUBREG\n")); err != nil {
		return err
	}
	select {
	case <-ds.okCh:
		return nil
	case errMsg := <-ds.errCh:
		return errors.New(errMsg.Payload.Error)
	}
}
