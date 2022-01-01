package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"manage/kaoshi/corejudge"
	"manage/kaoshi/loadfile"
)
type User struct{
	Id uint
	Username string//用户名
	Password string//密码
}
func UserLogin(c *gin.Context){//实现登录的控制器
	username := c.PostForm("username")
	password := c.PostForm("password")
	rows,err:=db.Query("select * from user")
	if err!=nil{
		fmt.Println("查询失败")
	}
	var s User
	for rows.Next(){
		err:=rows.Scan(&s.Id,&s.Username,&s.Password)
		if err !=nil{
			fmt.Println(err)
		}
	}
	if username!=s.Username{
		c.JSON(200,gin.H{
			"success":"false",
			"code":"404",
		})
	}else{
		//获取当前用户名的密码并查询是否匹配
		us,_:=db.Query("select password from user where username=?",s.Username)
		for us.Next(){
			var u User
			err :=us.Scan(&u.Password)
			if err!=nil{
				fmt.Println(err)
			}
			if password!=u.Password{
				c.JSON(200,gin.H{
					"success":"false",
					"code":404,
					"msg":"密码错误",
				})
			}else{
				c.JSON(200,gin.H{
					"success":"true",
					"code": "200",
					"msg":"登陆成功",
				})
			}
		}
	}
	rows.Close()
}
func UserRegister(c *gin.Context){//注册的函数
	username := c.PostForm("username")
	password := c.PostForm("password")
	rows,err:=db.Query("select * from user")
	if err!=nil{
		fmt.Println("查询失败")
	}
	var flag int64
	for rows.Next(){
		var s User
		err:=rows.Scan(&s.Username)
		if err != nil{
			fmt.Println(err)
		}
		flag=1
		if s.Username==username{
			flag=0
			break
		}
	}
	if flag==1 { //该用户名没有被注册
		_,err:=db.Exec("insert into user(username,password) values(?,?)",username,password)
		if err!=nil{
			fmt.Println(err)
			return
		}
	}else{//注册失败
		c.JSON(200,gin.H{
			"success":"false",
			"msg":"注册失败",
		})
	}
}
var db *sql.DB
func initDB(){//初始化数据库的函数
	dsn:="root:@tcp(127.0.0.1:3306)/student"
	db,_=sql.Open("mysql",dsn)
	db.SetConnMaxLifetime(10)
	db.SetMaxIdleConns(5)//设置最大闲置数
	if err:=db.Ping();err!=nil{
		fmt.Println("连接数据库失败")
		return
	}
	fmt.Println("连接数据库成功")
}
func LoadModel(c *gin.Context){
	c.HTML(200,"template.html",gin.H{})
}
func Add(c *gin.Context){//增加留言
	remain:=c.PostForm("string")
	num:=c.Query("id")//获取该登录用户的ID
	_,err:=db.Exec("update user set remain=? where id=?",remain,num)
	if err!=nil{
		fmt.Println(err)
		return
	}else{
		c.JSON(200,gin.H{
			"code":"200",
			"msg":"success",
		})
	}
}
func Delete(c *gin.Context){
	num:=c.Query("id")//获取该登录用户的ID
	_,err:=db.Exec("update user set remain=NULL where id=?",num)//将留言值为空
	if err!=nil{
		fmt.Println(err)
		return
	}else{
		c.JSON(200,gin.H{
			"code":"200",
			"msg":"success",
		})
	}
}
type To struct {
	remain string
}
func Search(c *gin.Context){//查询留言
	num:=c.Query("id")//获取用户id
	rows,err:=db.Query("select remain from user")
	if err!=nil{
		fmt.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next(){
		var u To
		err:=rows.Scan(&u.remain)
		if err!=nil{
			fmt.Println(err)
			return
		}
		c.JSON(200,gin.H{
			"ramin":u.remain,
			"id":num,
		})
	}
}
func StudentImport(c *gin.Context){
	score:=c.PostForm("score")
	id:=c.Query("id")
	//获取传来的ID和总分
	_,err:=db.Exec("update user set score=? where id=?",score,id)//插入分数中
	if err!=nil{
		fmt.Println(err)
		return
	}else{
		c.JSON(200,gin.H{
			"code":"200",
			"msg":"success",
		})
	}
}
type Score struct {
	score int64//分数
	id int64
}//存放分数的结构体
func StudentInput(c *gin.Context){//输出成绩
	rows,err:=db.Query("select score from user by score asc")//将查询的成绩按照升序排列
	if err!=nil{
		fmt.Println(err)
		return
	}
	var count=1
	defer rows.Close()
	for rows.Next(){
		var u Score
		err:=rows.Scan(&u.score,&u.id)
		if err!=nil{
			fmt.Println(err)
			return
		}
		c.JSON(200,gin.H{
			"score":u.score,
			"id":u.id,
			"scale":count,//该学生的分数
		})
		count++
	}
}
func main(){
	initDB()
	r:=gin.Default()
	r.LoadHTMLGlob("template/*")//加载模板
	r.GET("/lode",LoadModel)//打出模板
	user:=r.Group("/user")//登录注册
	{
		user.POST("/login",UserLogin)//登录的请求
		user.POST("/register",UserRegister)//注册的请求
	}
	file:=r.Group("/file")//文件的上传
	{
		file.POST("/convey",loadfile.FileConvey)
	}
	r.POST("/judge",corejudge.Judgenum)//实现签到功能
	remain:=r.Group("/remain")//实现留言的功能
	{
		remain.POST("/add",Add)//增加留言
		remain.DELETE("/delete",Delete)//删除留言
		remain.GET("/search",Search)//将留言显示出来
	}
	manage:=r.Group("/manage")//学生成绩管理
	{
		manage.POST("/student/import",StudentImport)//输入成绩
		manage.GET("student/output",StudentInput)//输出打印成绩
	}
	r.Run(":9090")//端口号
}
