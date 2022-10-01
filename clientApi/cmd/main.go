package main

import (
	"github.com/chatapp/clientapi/client"
	"github.com/chatapp/clientapi/internal/user"
	"github.com/chatapp/messagebroker"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	broker := messagebroker.NewRabbitMQ()
	defer broker.CloseConnection()

	userClient := client.NewUserClient()
	defer userClient.Close()

	publisher := messagebroker.NewRabbitPublish(broker)

	v1 := router.Group("/v1")
	{
		user.Route(v1, publisher, userClient)
	}

	router.Run("0.0.0.0:8080")
}
