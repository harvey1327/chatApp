package main

import (
	"github.com/gin-gonic/gin"
	"github.com/harvey1327/chatapp/clientapi/client"
	"github.com/harvey1327/chatapp/clientapi/internal/user"
	"github.com/harvey1327/chatapp/libmessagebroker"
)

func main() {
	router := gin.Default()
	broker := libmessagebroker.NewRabbitMQ()
	defer broker.CloseConnection()

	userClient := client.NewUserClient()
	defer userClient.Close()

	publisher := libmessagebroker.NewRabbitPublish(broker)

	v1 := router.Group("/v1")
	{
		user.Route(v1, publisher, userClient)
	}

	router.Run("0.0.0.0:8080")
}
