package handler

import (
	"manage/control"

	"github.com/gin-gonic/gin"
)
func Reboot(c *gin.Context){
	control.Reboot()
	c.JSON(200,gin.H{
		"code":200,
		"msg":"ok",
	})
}
func Shutdown(c *gin.Context){
	control.ShutDown()
	c.JSON(200,gin.H{
		"code":200,
		"msg":"ok",
	})
}