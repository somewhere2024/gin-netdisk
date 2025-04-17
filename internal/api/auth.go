package api

import (
	"fmt"
	"gin-netdisk/internal/models"
	"gin-netdisk/internal/schemas"
	"gin-netdisk/internal/services"
	"gin-netdisk/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// api/v1/user/login
func Login(c *gin.Context) {
	user := schemas.UserLogin{}

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "参数错误"})
		return
	}
	res, err := services.Authenticate(user.Username, user.Password)
	if err != nil {
		utils.Logger.Info("未授权")
		c.JSON(http.StatusOK, gin.H{"code": 403, "msg": "未授权", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "登录成功", "data": res})

}

// api/v1/user/register
func Register(c *gin.Context) {
	var user schemas.UserRegister
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "参数错误"})
		return
	}
	passwordHash, err := services.PasswordHash(user.Password)
	if err != nil {
		utils.Logger.Error("密码加密失败")
		return
	}
	newUser := models.User{
		Username:      user.Username,
		Password_hash: passwordHash,
		Email:         user.Email,
	}
	fmt.Println(newUser)
	rel := services.CreateUser(&newUser)
	if rel != nil {
		utils.Logger.Info("创建用户失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "注册成功", "data": nil})
}

// api/v1/user/profile
func UserProfile(c *gin.Context) { //需要保护路由

}
