package user

import (
	"net/http"

	"github.com/chatapp/messagebroker"
	"github.com/chatapp/messagebroker/events/createuser"
	"github.com/gin-gonic/gin"
)

func Route(parent *gin.RouterGroup, publisher messagebroker.Publish) {
	r := parent.Group("/user")
	{
		r.POST("/create", createUser(publisher))
	}
}

func createUser(publisher messagebroker.Publish) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request createuser.Model
		if err := ctx.BindJSON(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := publisher.Publish(messagebroker.PublishMessage(request, "user.create")); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, request)
	}
}
