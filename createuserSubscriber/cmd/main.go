package main

import (
	"log"

	"github.com/chatapp/messagebroker"
	"github.com/chatapp/messagebroker/events/createuser"
)

func main() {
	broker := messagebroker.NewRabbitMQ()
	defer broker.CloseConnection()

	broker.DeclareQueue("user.create")
	log.Println("listening on user.create")
	msgs := messagebroker.NewRabbitSubscribe[createuser.Model](broker).Subscribe(createuser.QUEUE_NAME)
	for {
		msg, ok := <-msgs
		if !ok {
			break
		}
		log.Printf("%+v", msg)
	}
	log.Println("END")
}
