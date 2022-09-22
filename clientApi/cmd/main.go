package main

import (
	"github.com/chatapp/clientapi/internal/user"
	"github.com/chatapp/messagebroker"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	broker := messagebroker.NewRabbitMQ()
	defer broker.CloseConnection()

	broker.DeclareQueue("user.create")
	commands := messagebroker.NewRabbitMQCommands[user.Request](broker)

	v1 := router.Group("/v1")
	{
		user.Route(v1, commands)
	}

	router.Run("0.0.0.0:8080")
}
