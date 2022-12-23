package room

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harvey1327/chatapp/clientapi/client"
	"github.com/harvey1327/chatapplib/messagebroker"
	"github.com/harvey1327/chatapplib/messagebroker/events/createroom"
)

func Route(parent *gin.RouterGroup, publisher messagebroker.Publish, roomClient client.RoomClient) {
	r := parent.Group("/room")
	{
		r.POST("/create", createRoom(publisher, roomClient))
	}
}

func createRoom(publisher messagebroker.Publish, roomClient client.RoomClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request createroom.Model
		if err := ctx.BindJSON(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		message := messagebroker.PublishMessage(request, createroom.QUEUE_NAME)
		if err := publisher.Publish(message); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result, err := roomClient.GetByEventID(message.EventMessage.EventID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, result)
	}
}
