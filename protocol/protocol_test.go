package protocol

import (
	"fmt"
	"testing"
)

func TestProtocolDecode(t *testing.T) {
	data := []byte(fmt.Sprintf("PUB %s\n", `{"topic":"default","message":"Hello World"}`))
	dp := new(DefaultProtocol)
	msg := &DefaultMessage{}
	dp.Decode(msg, data)
	if msg.Command != 2 || msg.Payload.Topic != "default" || msg.Payload.Message != "Hello World" {
		t.Errorf(`
		Error decoding:
		Expected: msg.Command=2, msg.Payload.Topic=default, msg.Payload.Message=Hello World
		Got: msg.Command=%d, msg.Payload.Topic=%s, msg.Payload.Message=%s`, msg.Command, msg.Payload.Topic, msg.Payload.Message)
	}
}

func TestMessageSerializer(t *testing.T) {
	rawPayload := []byte(`{"topic":"default","message":"Hello World"}`)
	dp := DefaultPayload{}
	err := dp.serialize(rawPayload)
	if err != nil {
		t.Errorf("Error serializing the payload")
	}

	if dp.Topic != "default" || dp.Message != "Hello World" {
		t.Errorf(`
		Expected topic=default and message=Hello World got topic=%s and message=%s`,
			dp.Topic, dp.Message)
	}
}
