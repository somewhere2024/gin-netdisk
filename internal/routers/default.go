package routers

import (
	"gin-netdisk/internal/api"
	"gin-netdisk/internal/middleware"
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
		// 跨域
		auth.Use(middleware.CURSMiddleware)

		auth.POST("/Login", api.Login)
		auth.POST("/Register", api.Register)
		auth.GET("/GetUserProfile", middleware.AuthMiddleware, api.GetUserProfile)
	}
	// 文件管理
	{
		file := r.Group("/api/v1/files")
		//file.Use(middleware.AuthMiddleware)
		file.POST("/upload", middleware.AuthMiddleware, api.FileUpload)
		file.GET("/:file_id/download", api.FileDownload)
		file.DELETE("/:file_id", middleware.AuthMiddleware, api.FileDelete)
		file.PUT("/:file_id", middleware.AuthMiddleware, api.FileRename) //重命名文件
		file.GET("/fileList", middleware.AuthMiddleware, api.FileGet)
	}
	//文件夹管理
	{
		folder := r.Group("/api/v1/folders")

		folder.POST("/:folder_id", middleware.AuthMiddleware, api.CreateFolder) //创建文件夹
		folder.GET("/:folder_id", middleware.AuthMiddleware, api.GetFolderList) //获取指定文件夹下的文件列表
	}

	r.Run("0.0.0.0:8000")
}
