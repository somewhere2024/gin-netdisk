package mysql

import (
	"fmt"
	"gin-netdisk/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {
	dsn := "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Print("连接数据库失败")
		log.Panic(err)
	}

}

func AutoMigrate() {
	err := DB.AutoMigrate(&models.User{}, &models.Resource{}, &models.Permission{}, &models.Trash{}, &models.Share{}) //自动迁移
	if err != nil {
		fmt.Print(err)
		panic("数据库迁移失败")
	}
}
