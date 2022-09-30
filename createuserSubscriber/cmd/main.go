package main

import (
	"log"

	"github.com/chatapp/database"
	"github.com/chatapp/messagebroker"
	"github.com/chatapp/messagebroker/events/createuser"
)

func main() {
	broker := messagebroker.NewRabbitMQ()
	defer broker.CloseConnection()
	broker.DeclareQueue(createuser.QUEUE_NAME)

	db := database.NewDB("user")
	defer db.Close()
	commands := database.NewCollection[messagebroker.EventMessage[createuser.Model]](db, "create")

	log.Println("listening on user.create")
	msgs := messagebroker.NewRabbitSubscribe[createuser.Model](broker).Subscribe(createuser.QUEUE_NAME)
	for {
		msg, ok := <-msgs
		if !ok {
			break
		}
		log.Printf("before %+v", msg)
		// save pending mesg to db
		err := commands.InsertOne(msg)
		if err != nil {
			log.Println(err)
		}
		// Check if userName exists, it will exist as we save the pending state
		existing, err := commands.FindSingleByQuery(database.Query("displayName", msg.Body.DisplayName))
		if err != nil {
			log.Println(err)
		}
		if existing.EventID != msg.EventID {
			//update msg to failed with error message and update db
			msg.Failed("username already exists")
		} else {
			msg.Complete()
		}

		log.Printf("after %+v", msg)
	}
	log.Println("END")
}
