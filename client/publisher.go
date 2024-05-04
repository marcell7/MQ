package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/marcell7/teenypubsub/protocol"
)

type Publisher interface {
	Publish() error  // Publishes user provided item/message to the specified topic
	start() error    // Starts listening for incoming messages
	connect() error  // Connects the client to the broker (tcp server)
	register() error // Registers the client as a publisher on the broker
}

// Implements Publisher interface
type DefaultPublisher struct {
	addr     string                        // Address of the broker
	conn     net.Conn                      // Connection of that specific subscriber - allows for writing and receiving messages from / to the broker
	protocol protocol.Protocol             // Protocol instance for decoding messages
	okCh     chan struct{}                 // Channel for signaling succesfully processed messages
	errCh    chan *protocol.DefaultMessage // Channel used for receiving messages from the broker
}

// Constructor for the DefaultPublisher struct
func NewPublisher(addr string) (*DefaultPublisher, error) {
	dp := &DefaultPublisher{
		addr:     addr,
		protocol: new(protocol.DefaultProtocol),
		okCh:     make(chan struct{}),
		errCh:    make(chan *protocol.DefaultMessage),
	}
	if err := dp.connect(); err != nil {
		return nil, err
	}

	go dp.start()

	if err := dp.register(); err != nil {
		return nil, err
	}
	return dp, nil
}

func (dp *DefaultPublisher) Publish(payload string) error {
	if _, err := dp.conn.Write([]byte(fmt.Sprintf("PUB %s\n", payload))); err != nil {
		return err
	}
	select {
	case <-dp.okCh:
		return nil
	case errMsg := <-dp.errCh:
		return errors.New(errMsg.Payload.Error)
	}
}

func (dp *DefaultPublisher) start() error {
	reader := bufio.NewReader(dp.conn)
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
		if err := dp.protocol.Decode(msg, data); err != nil {
			fmt.Println("Error decoding")
			return err
		}
		switch msg.Command {
		case protocol.CMD_OK:
			dp.okCh <- struct{}{}
		case protocol.CMD_ERROR:
			dp.errCh <- msg
		}
	}
}

func (dp *DefaultPublisher) connect() error {
	conn, err := net.Dial("tcp", dp.addr)
	if err != nil {
		return err
	}
	dp.conn = conn
	return nil
}

func (dp *DefaultPublisher) register() error {
	if _, err := dp.conn.Write([]byte("PUBREG\n")); err != nil {
		return err
	}
	select {
	case <-dp.okCh:
		return nil
	case errMsg := <-dp.errCh:
		return errors.New(errMsg.Payload.Error)
	}
}
