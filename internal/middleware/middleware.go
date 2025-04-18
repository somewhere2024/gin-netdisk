package middleware

import (
	"gin-netdisk/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// 解决跨域问题
func CURSMiddleware(c *gin.Context) {

}

// 授权中间件
func AuthMiddleware(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")

	if tokenStr == "" {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": 401, "msg": "未授权", "data": nil}) //终止请求的处理并返回
		return
	}
	token := strings.Split(tokenStr, " ")
	if len(token) != 2 {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": 401, "msg": "未授权", "data": nil}) //终止请求的处理并返回
		return
	}
	payload, err := services.GetCurrentUser(token[1])

	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": 401, "msg": "未授权", "data": nil}) //终止请求的处理并返回
		return
	}

	c.Set("userinfo", payload) //存储到上下文中

	c.Next()
}
