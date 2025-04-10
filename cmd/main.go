package main

import (
	"gin-netdisk/internal/app"
	"gin-netdisk/internal/config"
	"gin-netdisk/internal/utils"
	"log"
)

func main() {
	err := utils.InitLogger(config.Cfg)
	if err != nil {
		log.Print("初始化日志失败")
		log.Panic(err)
	}
	app.InitApp()
}
