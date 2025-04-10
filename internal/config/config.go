package config

type LogConfig struct {
	Level      string `json:"level"`       // 日志等级
	Filename   string `json:"filename"`    // 基准日志文件名
	MaxSize    int    `json:"maxsize"`     // 单个日志文件最大内容，单位：MB
	MaxAge     int    `json:"max_age"`     // 日志文件保存时间，单位：天
	MaxBackups int    `json:"max_backups"` // 最多保存几个日志文件
}

var Cfg = &LogConfig{
	Level:      "info",
	Filename:   "./logs/app.log",
	MaxSize:    100,
	MaxAge:     30,
	MaxBackups: 10,
}
