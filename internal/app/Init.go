package app

import (
	"gin-netdisk/internal/dao/mysql"
	"gin-netdisk/internal/routers"
)

func InitApp() {
	//初始化数据库
	mysql.InitDB()
	//sqlite再windows环境下需要cgo环境变量打开，并且还要gcc环境
	//sqlite.InitDB()

	//自动迁移
	mysql.AutoMigrate()
	
	routers.InitRouter()
}
