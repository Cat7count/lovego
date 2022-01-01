package loadfile

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

func FileConvey(c *gin.Context){//上传文件
	username:=c.PostForm("username")//获取传文件的人
	file,err:=c.FormFile("file")
	if err==nil{
		dst:=path.Join("./static/upload",file.Filename)//拼接上传文件的路径
		_=c.SaveUploadedFile(file,dst)
		c.JSON(200,gin.H{
			"success":"true",
			"username":username,
			"dst":"dst",
		})
	}
}
func FileLoad(c *gin.Context){//下载文件
	content := c.Query("content")
	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename=hello.zip")
	c.Header("Content-Type", "application/zip")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
	c.Writer.Write([]byte(content))
}