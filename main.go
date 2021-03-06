package main

import (
	"manage/handler"
	"manage/status"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
func main(){
	// 1. 获取要绑定的端口，如果没有->默认5200
	args := os.Args
	var port string
	if args == nil || len(args) < 2{
		port = ":5200"
	}else{
		port = args[1]
	}
	// 2. 绑定到对应端口

	go status.Update() //协程定时更新信息
	r := gin.Default()
	r.Use(Cors())
	status := r.Group("/status")
	{
		status.GET("/sys",handler.GetSysStatus)
		status.GET("/interval",handler.SetTimer)
		status.GET("/sysinfo",handler.GetSysInfo)
		status.GET("/check",handler.Check)
		status.GET("/all",handler.GetAllinfo)
	}
	control := r.Group("/control")
	{
		control.GET("/poweroff",handler.Shutdown)
		control.GET("/reboot",handler.Reboot)
	}
	r.Run(port)
}