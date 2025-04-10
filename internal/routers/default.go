package routers

import (
	"gin-netdisk/internal/services"
	"gin-netdisk/internal/utils"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.New()
	r.Use(utils.GinLogger(utils.Logger), utils.GinRecovery(utils.Logger, true))

	// 接口测试
	{
		rTest := r.Group("/test")
		rTest.GET("/status", services.TestStatus)

	}

	// 用户认证授权
	{
		r.Group("/auth")

	}

	// 文件管理
	{

	}
	r.Run("0.0.0.0:8000")
}
