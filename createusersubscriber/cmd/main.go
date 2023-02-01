package main

import (
	"log"

	"github.com/harvey1327/chatapp/createusersubscriber/config"
	"github.com/harvey1327/chatapplib/database"
	"github.com/harvey1327/chatapplib/messagebroker"
	"github.com/harvey1327/chatapplib/models/createuser"
)

func main() {
	conf := config.Load()
	broker := messagebroker.NewRabbitMQ(messagebroker.MessageBrokerConfig(conf.MB_HOST, conf.MB_PORT, conf.MB_USERNAME, conf.MB_PASSWORD))
	defer broker.CloseConnection()

	db := database.NewDB(database.USER, database.DBConfig(conf.DB_HOST, conf.DB_PORT, conf.DB_USERNAME, conf.DB_PASSWORD))
	defer db.Close()
	modelCol, eventCol := database.NewCollection[createuser.Model](db, createuser.GetModelConf().GetQueueName())

	msgs := messagebroker.NewRabbitSubscriber[createuser.Model](broker, createuser.GetModelConf().GetQueueName()).Subscribe()
	for {
		msg, ok := <-msgs
		if !ok {
			break
		}

		// save pending event mesg to db
		_, err := eventCol.InsertOne(msg)
		if err != nil {
			log.Panic(err)
		}
		// Check if userName exists in model collection
		_, err = modelCol.FindSingleByQuery(database.Query("displayName", msg.Body.DisplayName))
		if err != nil {
			if err == database.EMPTY {
				// Success that no existing username is found
				model, err := modelCol.InsertOne(msg.Body)
				if err != nil {
					log.Panic(err)
				}
				_, err = eventCol.InsertOne(msg.Complete(model.ID.Hex()))
				if err != nil {
					log.Panic(err)
				}
			} else {
				log.Panic(err)
			}
		} else {
			//If no error is thrown an existing userName already exists
			_, err = eventCol.InsertOne(msg.Failed("username already exists"))
			if err != nil {
				log.Panic(err)
			}
		}
	}
	log.Println("END")
}
