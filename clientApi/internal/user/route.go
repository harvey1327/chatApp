package user

import (
	"net/http"

	"github.com/chatapp/clientapi/internal/messagebroker"
	"github.com/gin-gonic/gin"
)

func Route(parent *gin.RouterGroup, service messagebroker.MessageBroker) {
	r := parent.Group("/user")
	{
		r.POST("/create", createUser(service))
	}
}

func createUser(service messagebroker.MessageBroker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request Request
		if err := ctx.BindJSON(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := service.Publish(messagebroker.Message(request, "user.create")); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, request)
	}
}
