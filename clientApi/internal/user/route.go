package user

import (
	"net/http"

	"github.com/chatapp/clientapi/client"
	"github.com/chatapp/libmessagebroker"
	"github.com/chatapp/libmessagebroker/events/createuser"
	"github.com/gin-gonic/gin"
)

func Route(parent *gin.RouterGroup, publisher libmessagebroker.Publish, userClient client.UserClient) {
	r := parent.Group("/user")
	{
		r.POST("/create", createUser(publisher, userClient))
	}
}

func createUser(publisher libmessagebroker.Publish, userClient client.UserClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request createuser.Model
		if err := ctx.BindJSON(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		message := libmessagebroker.PublishMessage(request, "user.create")
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
