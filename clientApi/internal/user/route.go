package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Route(parent *gin.RouterGroup) {
	r := parent.Group("/user")
	{
		r.POST("/create", createUser())
	}
}

func createUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request Request
		if err := ctx.BindJSON(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, request)
	}
}
