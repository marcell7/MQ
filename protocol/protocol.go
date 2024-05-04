package protocol

import (
	"bytes"
	"errors"
)

type Protocol interface {
	Decode(*DefaultMessage, []byte) error // Decodes a raw tcp message into -> {Command:<CMD_?>, Payload:<*DefaultPayload>}
}

type DefaultProtocol struct {
}

func (dp DefaultProtocol) Decode(msg *DefaultMessage, data []byte) error {
	var command []byte
	payload := new(DefaultPayload)
	args := bytes.Split(bytes.TrimSpace(data), []byte(" "))
	command = args[0]
	if len(args) >= 2 {
		pArgs := bytes.Join(args[1:], []byte(" "))
		if err := payload.serialize(pArgs); err != nil {
			return err
		}
	}
	switch string(command) {
	case "PUBREG":
		msg.Command = CMD_PUBREG
	case "SUBREG":
		msg.Command = CMD_SUBREG
	case "PUB":
		msg.Command = CMD_PUB
	case "SUB":
		msg.Command = CMD_SUB
	case "RECV":
		msg.Command = CMD_RECV
	case "RESP":
		msg.Command = CMD_RESP
	case "OK":
		msg.Command = CMD_OK
	case "ERROR":
		msg.Command = CMD_ERROR
	default:
		return errors.New("invalid command")
	}
	msg.Payload = payload
	return nil

}
