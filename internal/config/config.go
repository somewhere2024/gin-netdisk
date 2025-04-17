package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type LogConfig struct {
	Level      string `json:"level"`       // 日志等级
	Filename   string `json:"filename"`    // 基准日志文件名
	MaxSize    int    `json:"maxsize"`     // 单个日志文件最大内容，单位：MB
	MaxAge     int    `json:"max_age"`     // 日志文件保存时间，单位：天
	MaxBackups int    `json:"max_backups"` // 最多保存几个日志文件
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	driver   string
}

type Config struct {
	AppName   string
	AppPort   int
	DBConfig  DBConfig
	LogConfig LogConfig
}

func ConfigInit() {
	path, _ := os.Getwd()
	viper.AddConfigPath(path)
	viper.SetConfigName("config-env")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Fatalf("配置文件未找到: %v", err)
			} else {
				log.Fatalf("读取配置文件时出错: %v", err)
			}
		}
	}
}

func GetConfig() Config {
	ConfigInit()
	Cfg := Config{
		AppName: viper.GetString("app.name"),
		AppPort: viper.GetInt("app.port"),
		DBConfig: DBConfig{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetInt("database.port"),
			User:     viper.GetString("database.username"),
			Password: viper.GetString("database.password"),
			Name:     viper.GetString("database.name"),
			driver:   viper.GetString("database.driver"),
		},
		LogConfig: LogConfig{
			Level:      viper.GetString("log.level"),
			Filename:   viper.GetString("log.filePath"),
			MaxSize:    viper.GetInt("log.maxSize"),
			MaxAge:     viper.GetInt("log.maxAge"),
			MaxBackups: viper.GetInt("log.maxBackups"),
		},
	}
	return Cfg
}

var Cfg = GetConfig()
