package app

import (
	"gin-netdisk/internal/dao/mysql"
	"gin-netdisk/internal/routers"
)

func InitApp() {
	//初始化数据库
	mysql.InitDB()
	routers.InitRouter()

}
