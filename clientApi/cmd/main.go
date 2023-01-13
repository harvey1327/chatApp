package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/harvey1327/chatapp/clientapi/client"
	"github.com/harvey1327/chatapp/clientapi/config"
	"github.com/harvey1327/chatapp/clientapi/internal/room"
	"github.com/harvey1327/chatapp/clientapi/internal/user"
	"github.com/harvey1327/chatapplib/messagebroker"
	"github.com/harvey1327/chatapplib/models/createroom"
	"github.com/harvey1327/chatapplib/models/createuser"
)

func main() {
	conf := config.Load()
	router := gin.Default()
	broker := messagebroker.NewRabbitMQ(messagebroker.MessageBrokerConfig(conf.MB_HOST, conf.MB_PORT, conf.MB_USERNAME, conf.MB_PASSWORD))
	defer broker.CloseConnection()

	userClient := client.NewUserClient(conf.USER_SERVICE_HOST, conf.USER_SERVICE_PORT)
	defer userClient.Close()

	roomClient := client.NewRoomClient(conf.ROOM_SERVICE_HOST, conf.ROOM_SERVICE_PORT)
	defer roomClient.Close()

	createUserPublisher := messagebroker.NewRabbitPublisher[createuser.Model](broker, createuser.GetModelConf().GetQueueName())
	createRoomPublisher := messagebroker.NewRabbitPublisher[createroom.Model](broker, createroom.GetModelConf().GetQueueName())

	v1 := router.Group("/v1")
	{
		user.Route(v1, createUserPublisher, userClient)
		room.Route(v1, createRoomPublisher, roomClient)
	}

	router.Run(fmt.Sprintf("%s:%d", conf.HOST, conf.PORT))
}
