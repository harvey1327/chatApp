package room

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harvey1327/chatapp/clientapi/client"
	"github.com/harvey1327/chatapplib/messagebroker"
	"github.com/harvey1327/chatapplib/models/createroom"
)

func Route(parent *gin.RouterGroup, publisher messagebroker.Publish[createroom.Model], roomClient client.RoomClient) {
	r := parent.Group("/room")
	{
		r.POST("/create", createRoom(publisher, roomClient))
	}
}

func createRoom(publisher messagebroker.Publish[createroom.Model], roomClient client.RoomClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request createroom.Model
		if err := ctx.BindJSON(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if message, err := publisher.Publish(request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			events, err := roomClient.GetByEventID(message.EventID)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			var response any
			for event := range events {
				log.Printf("Read event: %+v\n", event)
				response = event
			}
			ctx.JSON(200, response)
		}
	}
}
