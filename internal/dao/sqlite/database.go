package sqlite

import (
	"fmt"
	"gin-netdisk/internal/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error

	// 使用绝对路径或配置文件来设置数据库文件路径
	// 如果希望使用相对路径，可以调整路径
	dbPath := "test.db"
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		// 使用 fmt.Errorf 提供更多的上下文信息
		utils.Logger.Panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
}
