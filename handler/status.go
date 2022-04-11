package handler

import (
	"manage/status"
	"strconv"

	"github.com/gin-gonic/gin"
)
func Check(c *gin.Context){
	c.JSON(200,gin.H{
		"code":200,
		"msg":"ok",
	})
}

func GetSysStatus(c *gin.Context){
	c.JSON(200,gin.H{
		"code":200,
		"data":status.GetSysDataAddr(),
	})
}

func GetSysInfo(c *gin.Context){

	c.JSON(200,gin.H{
		"code":200,
		"data":status.GetSysInfo(),
	})
}
func GetAllinfo(c *gin.Context){
	c.JSON(200,gin.H{
		"code":200,
		"data":status.GetTotal(),
	})
}

func SetTimer(c *gin.Context){
	time := c.Query("time")
	if time != ""{
		time_i,_ := strconv.ParseInt(time,10,64)
		status.SetInterVal(time_i)
		c.JSON(200,gin.H{
			"code":200,
			"msg":"set time:" + time,
		})
		return
	}
	c.JSON(200,gin.H{
		"code":500,
		"msg":"please set interval",
	})
}
