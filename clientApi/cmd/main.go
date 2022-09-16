package main

import (
	"github.com/chatapp/clientapi/internal/messagebroker"
	"github.com/chatapp/clientapi/internal/user"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	broker := messagebroker.NewRabbitMQ()
	defer broker.CloseConnection()

	broker.DeclareQueue("user.create")

	v1 := router.Group("/v1")
	{
		user.Route(v1, broker)
	}

	router.Run("0.0.0.0:8080")
}
