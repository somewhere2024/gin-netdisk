package routers

import (
	"gin-netdisk/internal/api"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	//r := gin.New()
	////使用自定义日志中间件
	//r.Use(utils.GinLogger(utils.Logger), utils.GinRecovery(utils.Logger, true))
	r := gin.Default()

	// 接口测试
	{
		rTest := r.Group("/api/v1/test")
		rTest.GET("/status", api.TestStatus)

	}

	// 用户认证授权
	{
		auth := r.Group("/api/v1/auth")
		auth.POST("/Login", api.Login)
		auth.POST("/Register", api.Register)
		auth.GET("/GetUserProfile")

	}

	// 文件管理
	{

	}

	r.Run("0.0.0.0:8000")
}
