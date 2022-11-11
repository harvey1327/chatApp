package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harvey1327/chatapp/clientapi/client"
	"github.com/harvey1327/chatapplib/messagebroker"
	"github.com/harvey1327/chatapplib/messagebroker/events/createuser"
)

func Route(parent *gin.RouterGroup, publisher messagebroker.Publish, userClient client.UserClient) {
	r := parent.Group("/user")
	{
		r.POST("/create", createUser(publisher, userClient))
	}
}

func createUser(publisher messagebroker.Publish, userClient client.UserClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request createuser.Model
		if err := ctx.BindJSON(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		message := messagebroker.PublishMessage(request, createuser.QUEUE_NAME)
		if err := publisher.Publish(message); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result, err := userClient.GetByEventID(message.EventMessage.EventID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, result)
	}
}
