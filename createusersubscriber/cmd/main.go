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

		// save pending mesg to db
		insert, err := commands.InsertOne(msg)
		if err != nil {
			log.Println(err)
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
			log.Println(err)
		}
	}
	log.Println("END")
}
