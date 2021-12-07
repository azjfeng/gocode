package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")
		context.Next()
	}
}

func main() {
	router := gin.Default()
	router.Use(Cors())

	router.POST("/post", func(c *gin.Context) {
		// c.Header("Access-Control-Allow-Origin", "*")
		id := c.DefaultQuery("id", "11")
		page := c.DefaultQuery("page", "0")

		fmt.Printf("id: %s; page: %s", id, page)
		c.JSON(200, gin.H{
			"name": "jameinfeng",
		})
	})
	router.Run(":8080")
}
