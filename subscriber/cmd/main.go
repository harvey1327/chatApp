package main

import (
	"log"

	"github.com/chatapp/messagebroker"
	"github.com/chatapp/subscriber/internal/user"
)

func main() {
	broker := messagebroker.NewRabbitMQ()
	defer broker.CloseConnection()

	broker.DeclareQueue("user.create")
	log.Println("listening on user.create")
	msgs := messagebroker.NewRabbitMQCommands[user.CreateUser](broker).Subscribe("user.create")
	for {
		msg, ok := <-msgs
		if !ok {
			break
		}
		log.Printf("%+v", msg)
	}
	log.Println("END")
}
