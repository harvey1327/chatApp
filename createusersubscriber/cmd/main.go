package main

import (
	"log"

	"github.com/harvey1327/chatapplib/database"
	"github.com/harvey1327/chatapplib/messagebroker"
	"github.com/harvey1327/chatapplib/messagebroker/events/createuser"
)

func main() {
	broker := messagebroker.NewRabbitMQ()
	defer broker.CloseConnection()

	db := database.NewDB(database.USER)
	defer db.Close()
	commands := database.NewCollection[messagebroker.EventMessage[createuser.Model]](db, createuser.QUEUE_NAME)

	log.Printf("listening on %s\n", createuser.QUEUE_NAME)
	msgs := messagebroker.NewRabbitSubscribe[createuser.Model](broker).Subscribe(createuser.QUEUE_NAME)
	for {
		msg, ok := <-msgs
		if !ok {
			break
		}

		// save pending mesg to db
		insert, err := commands.InsertOne(msg)
		if err != nil {
			log.Panic(err)
		}
		// Check if userName exists, it will exist as we save the pending state
		existing, err := commands.FindSingleByQuery(database.Query("displayName", insert.Data.Body.DisplayName))
		if err != nil {
			insert.Data.Failed(err.Error())
		}
		if existing.Data.EventID != insert.Data.EventID {
			//update msg to failed with error message and update db
			insert.Data.Failed("username already exists")
		} else {
			insert.Data.Complete()
		}

		//save to db
		err = commands.FindByIDAndUpdate(insert)
		if err != nil {
			log.Panic(err)
		}
	}
	log.Println("END")
}
