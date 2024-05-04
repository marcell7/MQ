package main

import (
	"log"

	"github.com/marcell7/teenypubsub/broker"
)

func main() {
	b := broker.New(":3000")
	b.AddTopic("default")
	go b.Listen()
	log.Println("Broker is ready. Listening on port 3000")
	select {}
}
