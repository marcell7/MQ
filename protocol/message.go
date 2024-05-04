package protocol

import "encoding/json"

type Message interface {
}

type DefaultMessage struct {
	Id      string
	Command Command
	Payload *DefaultPayload
}

type Payload interface {
	serialize([]byte) error
}

type DefaultPayload struct {
	Topic   string
	Message string
	Error   string
}

func (dp *DefaultPayload) serialize(rawPayload []byte) error {
	err := json.Unmarshal(rawPayload, dp)
	if err != nil {
		return err
	}
	return nil
}

type DefaultResponsePayload struct {
	Data []byte
}

func (rp *DefaultResponsePayload) serialize(rawPayload []byte) error {
	err := json.Unmarshal(rawPayload, rp)
	if err != nil {
		return err
	}
	return nil
}
