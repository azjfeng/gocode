package main

import (
	"fmt"
	"io"
	"log"
	"os"

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

type Person struct {
	User     string `form:"user" binding:"required"`
	PassWord string `form:"password" binding:"required"`
}

func main() {
	router := gin.Default()

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	// 如果需要同时将日志写入文件和控制台，请使用以下代码。
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router.Use(Cors())
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Static("/", "./public")
	router.POST("/post", func(c *gin.Context) {
		// c.Header("Access-Control-Allow-Origin", "*")
		id := c.DefaultQuery("id", "11")
		page := c.DefaultQuery("page", "0")

		fmt.Printf("id: %s; page: %s", id, page)
		c.JSON(200, gin.H{
			"name": "jameinfeng",
		})
	})

	router.POST("/upload", func(c *gin.Context) {
		/** 单文件上传 */
		// file, _ := c.FormFile("file")
		// log.Println(file)
		// c.SaveUploadedFile(file, "./public/"+file.Filename)

		/** 多文件上传 */
		form, _ := c.MultipartForm()
		files := form.File["file"]

		for _, file := range files {
			log.Println(file.Filename)

			// 上传文件至指定目录
			c.SaveUploadedFile(file, "./public/"+file.Filename)
		}

		c.JSON(200, gin.H{
			"message": "上传成功",
		})
	})

	authorized := router.Group("/cgi")
	{
		authorized.POST("/login", func(c *gin.Context) { c.String(200, "1") })
		authorized.POST("/submit", func(c *gin.Context) { c.String(200, "1") })
		authorized.POST("/read", func(c *gin.Context) { c.String(200, "1") })

		// 嵌套路由组
		testing := authorized.Group("testing")
		testing.POST("/analytics", func(c *gin.Context) { c.String(200, "1") })
	}

	router.POST("/getUserInfo", func(c *gin.Context) {
		var person Person
		if (c.ShouldBindJSON(&person)) == nil {
			log.Println(person.User)
			log.Println(person.PassWord)
			c.JSON(200, gin.H{
				"name": "jamefeine",
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "密码错误",
		})
	})
	router.Run(":8080")
}
