package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID     int
	Auther string
	Title  string
}

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
	db := InitDB()
	// var user User
	user := User{Auther: "Jinzhu", Title: "测试"}
	db.Create(&user)
	fmt.Println(user)
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	// 如果需要同时将日志写入文件和控制台，请使用以下代码。
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router.Use(Cors())
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	// router.Static("/", "./public")
	authorized := router.Group("/cgi")
	{
		authorized.POST("/login", func(c *gin.Context) { c.String(200, "1") })
		authorized.POST("/submit", func(c *gin.Context) { c.String(200, "1") })
		authorized.POST("/read", func(c *gin.Context) { c.String(200, "1") })
		authorized.POST("/upload", func(c *gin.Context) {
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
				c.SaveUploadedFile(file, "/usr/local/static/"+file.Filename)
			}

			c.JSON(200, gin.H{
				"message": "上传成功",
			})
		})
		authorized.POST("/getUserInfo", func(c *gin.Context) {
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

		// 嵌套路由组
		testing := authorized.Group("testing")
		testing.POST("/analytics", func(c *gin.Context) { c.String(200, "1") })
	}

	router.Run(":12345")
}

func InitDB() *gorm.DB {
	//前提是你要先在本机用Navicat创建一个名为go_db的数据库
	host := "localhost"
	port := "3306"
	database := "reactBoke"
	username := "root"
	password := "123456"
	charset := "utf8"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)
	//这里 gorm.Open()函数与之前版本的不一样，大家注意查看官方最新gorm版本的用法
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error to Db connection, err: " + err.Error())
	}
	//这个是gorm自动创建数据表的函数。它会自动在数据库中创建一个名为users的数据表
	_ = db.AutoMigrate(&User{})
	return db
}
