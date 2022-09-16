package main

import (
	"log"

	"github.com/chatapp/subscriber/internal/messagebroker"
)

func main() {
	broker := messagebroker.NewRabbitMQ()
	defer broker.CloseConnection()

	broker.DeclareQueue("user.create")
	log.Println("listening on user.create")
	msgs := broker.Subscribe("user.create")
	for {
		msg, ok := <-msgs
		if !ok {
			break
		}
		log.Printf("%+v", msg)
	}
	log.Println("END")
}
