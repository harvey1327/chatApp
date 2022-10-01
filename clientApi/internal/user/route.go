package user

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chatapp/messagebroker"
	"github.com/chatapp/messagebroker/events/createuser"
	"github.com/chatapp/proto/generated/userpb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		message := messagebroker.PublishMessage(request, "user.create")
		if err := publisher.Publish(message); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "0.0.0.0", 50052), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal(err)
		}
		userClient := userpb.NewServiceClient(conn)
		time.Sleep(3 * time.Second)
		result, err := userClient.GetByEventID(context.TODO(), &userpb.GetByEventIDRequest{EventID: message.EventMessage.EventID})
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, result)
	}
}
