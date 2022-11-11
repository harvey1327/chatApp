package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/harvey1327/chatapp/clientapi/client"
	"github.com/harvey1327/chatapp/clientapi/config"
	"github.com/harvey1327/chatapp/clientapi/internal/user"
	"github.com/harvey1327/chatapplib/messagebroker"
)

func main() {
	conf := config.Load()
	router := gin.Default()
	broker := messagebroker.NewRabbitMQ(messagebroker.MessageBrokerConfig(conf.MB_HOST, conf.MB_PORT, conf.MB_USERNAME, conf.MB_PASSWORD))
	defer broker.CloseConnection()

	userClient := client.NewUserClient(conf.USER_SERVICE_HOST, conf.USER_SERVICE_PORT)
	defer userClient.Close()

	publisher := messagebroker.NewRabbitPublish(broker)

	v1 := router.Group("/v1")
	{
		user.Route(v1, publisher, userClient)
	}

	router.Run(fmt.Sprintf("%s:%d", conf.HOST, conf.PORT))
}
