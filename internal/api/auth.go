package api

import (
	"gin-netdisk/internal/dao/mysql"
	"gin-netdisk/internal/models"
	"gin-netdisk/internal/schemas"
	"gin-netdisk/internal/services"
	"gin-netdisk/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	toEncode := jwt.MapClaims{"user_id": res.ID}

	token, err := services.CreateToken(toEncode)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "登录成功", "access_token": token, "token_type": "bearer", "data": nil})

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
	rel := services.CreateUser(&newUser)
	if rel != nil {
		utils.Logger.Info("创建用户失败")
		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "创建用户失败", "data": nil})
		return
	}

	// 创建用户对应的文件夹
	err = services.UserCreateFolder(newUser.ID)
	if err != nil {
		utils.Logger.Error("创建用户文件夹失败")
		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "创建用户文件夹失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "注册成功", "data": nil})
}

// api/v1/user/profile
func GetUserProfile(c *gin.Context) {
	userInfo, exists := c.Get("userinfo")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "未授权", "data": nil})
		return
	}

	user_id := userInfo.(jwt.MapClaims)["user_id"]

	user := models.User{}

	rel := mysql.DB.Where("id = ?", user_id).First(&user)
	if rel.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "未授权", "data": nil})
		return
	}

	userRes := models.UserResponse{
		Username:     user.Username,
		CreatedAt:    user.CreatedAt,
		Email:        user.Email,
		StorageTotal: user.StorageTotal,
		StorageUsed:  user.StorageUsed,
		UpdateAt:     user.UpdatedAt,
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取用户信息成功", "data": userRes})
}
