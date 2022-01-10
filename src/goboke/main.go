package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Db *sqlx.DB

type Info struct {
	Id    int    `Db:"id"`
	Age   int    `Db:"age"`
	Sex   int    `Db:"sex"`
	Name  string `Db:"name"`
	Phone string `Db:"phone"`
}

type ShareList struct {
	Id          int    `Db:"id"`
	Auther      string `Db:"auther"`
	Title       string `Db:"title"`
	Create_Time string `Db:"create_time"`
	Content     string `Db:"content"`
	Support     int    `Db:"support"`
	Watch_Num   int    `Db:"watch_num"`
	Image       string `Db:"image"`
	Contentdesc string ` Db:"contentdesc"`
}

const (
	user     = "root"
	password = "123456"
	host     = "127.0.0.1:3306"
	dbname   = "reactboke"
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
	initDB()

	//writeFile()
	authorized := router.Group("/cgi")
	{
		authorized.POST("/login", func(c *gin.Context) {
			info := []Info{}
			Db.Select(&info, "select * from user")
			c.JSON(200, info)
		})
		authorized.POST("/submit", func(c *gin.Context) { c.String(200, "1") })
		authorized.POST("/read", func(c *gin.Context) { c.String(200, "1") })

		authorized.POST("/getTechnologyShare", func(c *gin.Context) {
			sharelist := []ShareList{}
			err := Db.Select(&sharelist, "select * from technology_share")
			fmt.Println(err)
			c.JSON(200, sharelist)
		})

		authorized.POST("/addTechnologyShare", func(c *gin.Context) {

			// 24小时制
			timeObj := time.Now()
			var str = timeObj.Format("2006/01/02 15:04:05")
			fmt.Println(str) // 2020/04/26 17:48:53

			tx := Db.MustBegin()
			err := tx.MustExec("insert into technology_share (id,auther,title,create_time,content,support,watch_num,image,contentdesc) values ($1,$2,$3,$4,$5,$6,$7,$8,$9)", "null", "jamefeng", "test", str, "", 1, 1, "https://www.azjfeng.com/0cdba396-4569-47aa-ae61-ed788dbf6f84.jpg")
			tx.Commit()
			fmt.Println(err)
			c.JSON(200, gin.H{"message": "添加成功"})
		})
	}

	router.Run(":3332")
}
func initDB() {
	//数据库连接
	db, _ := sqlx.Open("mysql", user+":"+password+"@tcp("+host+")/"+dbname)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	Db = db
	err := db.Ping()
	if err != nil {
		fmt.Println("数据库链接失败")
	}
	////多行查询
	//rows,e:=db.Query("select * from user")
	//
	//fmt.Println(e)
	//fmt.Println("rows",rows)
	//var id ,age ,sex int ; var name, phone string
	//for rows.Next(){
	//	err :=rows.Scan(&id,&name, &age, &sex, &phone)
	//	if err != nil {
	//		fmt.Println("get data failed, error:[%v]", err.Error())
	//	}
	//	fmt.Println(id,name,age,sex,phone)
	//}
	//defer db.Close()
}

//写文件
func writeFile() {
	userFile := "D://test.txt"
	f, err := os.Create(userFile)
	if err != nil {
		fmt.Println(userFile, err)
		return
	}
	defer f.Close()
	for i := 0; i < 10; i++ {
		f.WriteString("www.361way.com22232311222!\n")
		f.Write([]byte("Just a test!\n"))
	}
}
