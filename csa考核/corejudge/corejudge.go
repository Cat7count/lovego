package corejudge

import "github.com/gin-gonic/gin"

func Judgenum(c *gin.Context){
	judgenum:=c.PostForm("judgenum")//获取前端传来的数字
	if judgenum == "1"{
		c.JSON(200,gin.H{
			"success":"false",
			"msg":"签到成功",
		})
	}else{
		c.JSON(200,gin.H{
			"success":"false",
			"msg":"签到失败",
		})
	}
}