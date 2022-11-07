package main

import (
	"log"

	"github.com/chatapp/libdatabase"
	"github.com/chatapp/libmessagebroker"
	"github.com/chatapp/libmessagebroker/events/createuser"
)

func main() {
	broker := libmessagebroker.NewRabbitMQ()
	defer broker.CloseConnection()

	db := libdatabase.NewDB("user")
	defer db.Close()
	commands := libdatabase.NewCollection[libmessagebroker.EventMessage[createuser.Model]](db, "create")

	log.Println("listening on user.create")
	msgs := libmessagebroker.NewRabbitSubscribe[createuser.Model](broker).Subscribe(createuser.QUEUE_NAME)
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
		existing, err := commands.FindSingleByQuery(libdatabase.Query("displayName", insert.Data.Body.DisplayName))
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
