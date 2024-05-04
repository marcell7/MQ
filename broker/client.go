package broker

import (
	"fmt"
	"net"
)

type Client interface {
	sendResp(*Item) error
	sendOk() error
	sendError(string) error
}

type Publisher struct {
	id   string
	conn net.Conn
}

func newPublisher(id string, conn net.Conn) *Publisher {
	return &Publisher{
		id:   id,
		conn: conn,
	}
}

func (p *Publisher) sendOk() error {
	_, err := p.conn.Write([]byte("OK\n"))
	if err != nil {
		return err
	}
	return nil
}

func (p *Publisher) sendResp(item *Item) error {
	return nil
}

func (p *Publisher) sendError(errorMsg string) error {
	payload := fmt.Sprintf(`{"error":"%s"}`, errorMsg)
	_, err := p.conn.Write([]byte(fmt.Sprintf("ERROR %s\n", payload)))
	if err != nil {
		return err
	}
	return nil
}

type Subscriber struct {
	id   string
	conn net.Conn
}

func newSubscriber(id string, conn net.Conn) *Subscriber {
	return &Subscriber{
		id:   id,
		conn: conn,
	}
}

func (s *Subscriber) sendOk() error {
	_, err := s.conn.Write([]byte("OK\n"))
	if err != nil {
		return err
	}
	return nil
}

func (s *Subscriber) sendResp(item *Item) error {
	payload := fmt.Sprintf(`{"message":"%s"}`, item.Data)
	_, err := s.conn.Write([]byte(fmt.Sprintf("RESP %s\n", payload)))
	if err != nil {
		return err
	}
	return nil

}

func (s *Subscriber) sendError(errorMsg string) error {
	payload := fmt.Sprintf(`{"error":"%s"}`, errorMsg)
	_, err := s.conn.Write([]byte(fmt.Sprintf("ERROR %s\n", payload)))
	if err != nil {
		return err
	}
	return nil
}
