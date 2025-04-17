package utils

import (
	"gin-netdisk/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// Logger 是全局的日志实例
var Logger *zap.Logger

// InitLogger 初始化日志记录器，基于配置文件设置日志级别、文件名等参数
// 参数：
//   - cfg: 日志配置结构体指针，包含日志文件名、最大大小、备份数量、过期天数和日志级别
//
// 返回值：
//   - err: 初始化过程中可能发生的错误
func InitLogger(cfg *config.Config) (err error) {
	// 创建文件写入器
	fileWriteSyncer := getLogWriter(cfg.LogConfig.Filename, cfg.LogConfig.MaxSize, cfg.LogConfig.MaxBackups, cfg.LogConfig.MaxAge)
	// 创建标准输出写入器
	consoleWriteSyncer := zapcore.AddSync(os.Stdout)

	// 创建一个多路写入器，将日志同时输出到文件和标准输出
	writeSyncer := zapcore.NewMultiWriteSyncer(fileWriteSyncer, consoleWriteSyncer)

	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(cfg.LogConfig.Level))
	if err != nil {
		return
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)

	Logger = zap.New(core, zap.AddCaller())
	return
}

// getEncoder 配置并返回一个 zap 编码器
// 该函数定义了日志的时间格式、日志级别格式、持续时间格式以及调用者信息的编码方式
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// getLogWriter 返回一个使用 lumberjack 的日志写入器
// 该函数实现了日志文件的轮转逻辑，包括文件大小限制、备份文件数量和日志保留天数
func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// GinLogger 记录 gin 框架的默认日志
// 参数：
//   - logger: zap 日志实例
//
// 功能：
//   - 记录每个请求的详细信息，包括状态码、方法、路径、查询参数、IP 地址、用户代理、错误信息和耗时
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery 捕获项目中可能出现的 panic，并使用 zap 记录相关日志
// 参数：
//   - logger: zap 日志实例
//   - stack: 是否记录堆栈信息
//
// 功能：
//   - 捕获 panic 并根据情况记录错误日志，支持是否记录堆栈信息
//   - 如果连接中断（如 broken pipe），则不会记录完整的堆栈信息
func GinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 检查是否为连接中断错误（broken pipe 或 connection reset by peer）
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// 如果连接已断开，则无法向其写入状态
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
