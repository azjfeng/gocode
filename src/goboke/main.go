package main

import (
	"fmt"
	"io/ioutil"
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
	Id          int    `json:"id"`
	Auther      string `json:"auther"`
	Title       string `json:"title"`
	Create_Time string `json:"create_time"`
	Content     string `json:"content"`
	Support     int    `json:"support"`
	Watch_Num   int    `json:"watch_num"`
	Image       string `json:"image"`
	Contentdesc string ` json:"contentdesc"`
}

type Form struct {
	Title string `json:"title"`
}

type AddData struct {
	Title string `json:"title"`
	Auther string `json:"auther"`
	Desc string `json:"desc"`
	Content string `json:"content"`
}

type UpdateData struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Auther string `json:"auther"`
	Desc string `json:"desc"`
	Content string `json:"content"`
}

const (
	user     = "root"
	password = "123456"
	host     = "127.0.0.1:3306"
	dbname   = "reactBoke"
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
	authorized := router.Group("/common")
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
			// 声明接收的变量
			var json AddData

			// 将request的body中的数据，自动按照json格式解析到结构体
			if err := c.ShouldBindJSON(&json); err != nil {
				// 返回错误信息
				// gin.H封装了生成json数据的工具
				c.JSON(-1, gin.H{"error": err.Error()})
				return
			}
			writeFile("/usr/local/static/text/" + json.Title + ".txt", json.Content)
			fmt.Println(json)
			// 24小时制
			timeObj := time.Now()
			var str = timeObj.Format("2006/01/02 15:04:05")
			fmt.Println(str) // 2020/04/26 17:48:53

			tx := Db.MustBegin()
			err := tx.MustExec("insert into technology_share (auther,title,create_time,content,support,watch_num,image,contentdesc) values (?,?,?,?,?,?,?,?)", json.Auther, json.Title, str, "", 1, 1, "https://www.azjfeng.com/static/2019-06-20-1.png",json.Desc)
			tx.Commit()
			fmt.Println(err)
			c.JSON(200, gin.H{"message": "添加成功"})
		})

		authorized.POST("/getDetail", func(c *gin.Context) {

			// 声明接收的变量
			var json Form
			// 将request的body中的数据，自动按照json格式解析到结构体
			if err := c.ShouldBindJSON(&json); err != nil {
				// 返回错误信息
				// gin.H封装了生成json数据的工具
				c.JSON(-1, gin.H{"error": err.Error()})
				return
			}
			content, _ := ReadAll("/usr/local/static/text/" + json.Title + ".txt")
			fmt.Println(string(content))
			sharelist := []ShareList{}
			err := Db.Select(&sharelist, "select * from technology_share where title = ?", json.Title)
			fmt.Println(err)
			c.JSON(200, gin.H{"result": sharelist, "content": string(content)})
		})

		authorized.POST("/updateTechnologyShare", func(c *gin.Context) {
			//// 声明接收的变量
			var json UpdateData

			// 将request的body中的数据，自动按照json格式解析到结构体
			if err := c.ShouldBindJSON(&json); err != nil {
				// 返回错误信息
				// gin.H封装了生成json数据的工具
				c.JSON(-1, gin.H{"error": err.Error()})
				return
			}
			writeFile("/usr/local/static/text/" + json.Title + ".txt", json.Content)
			sqlStr := "update technology_share set title = ?,  auther= ?, contentdesc=? where id = ?"
			_, err := Db.Exec(sqlStr, json.Title, json.Auther, json.Desc,json.Id)
			if err != nil {
				fmt.Printf("update failed, err:%v\n", err)
				return
			}
			c.JSON(200, gin.H{"message": "修改成功"})
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
func writeFile(path string, content string) {
	userFile := path
	f, err := os.Create(userFile)
	if err != nil {
		fmt.Println(userFile, err)
		return
	}
	defer f.Close()
	f.WriteString(content)
}

func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}
