package main

import (
	"github.com/chatapp/clientapi/internal/user"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		user.Route(v1)
	}

	router.Run(":8080")
}
