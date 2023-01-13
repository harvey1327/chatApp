package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harvey1327/chatapp/clientapi/client"
	"github.com/harvey1327/chatapplib/messagebroker"
	"github.com/harvey1327/chatapplib/models/createuser"
)

func Route(parent *gin.RouterGroup, publisher messagebroker.Publish[createuser.Model], userClient client.UserClient) {
	r := parent.Group("/user")
	{
		r.POST("/create", createUser(publisher, userClient))
	}
}

func createUser(publisher messagebroker.Publish[createuser.Model], userClient client.UserClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request createuser.Model
		if err := ctx.BindJSON(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if message, err := publisher.Publish(request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			result, err := userClient.GetByEventID(message.EventID)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(200, result)
		}
	}
}
