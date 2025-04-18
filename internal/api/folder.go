package api

import (
	"fmt"
	"gin-netdisk/internal/dao/mysql"
	"gin-netdisk/internal/models"
	"gin-netdisk/internal/schemas"
	"gin-netdisk/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

/*
创建文件夹
通过路径参数获取父文件夹id，通过表单获取文件夹名
如果父文件夹id为空，则创建在用户根目录下
*/
func CreateFolder(c *gin.Context) {
	parentId := c.Param("folder_id")
	folderName := c.PostForm("folder_name")
	userInfo, _ := c.Get("userinfo")
	fmt.Println(parentId, folderName)
	userId := userInfo.(jwt.MapClaims)["user_id"].(string)
	if len(parentId) < 36 {
		if err := mysql.DB.Model(&models.Resource{}).Where("user_id = ? and name = ?", userId, userId).Select("id").Scan(&parentId).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "创建文件夹失败", "data": nil})
			return
		}
	}
	path := services.GetFolderParentPath(parentId) + "/" + folderName

	fmt.Println(folderName, parentId)
	if err := services.CreateFolder(userId, folderName, parentId, path); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "创建文件夹失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建文件夹成功", "data": nil})
}

/*
根据用户id和文件夹id获取该文件夹下的文件夹列表

查询参数folder_id

默认会获取用户根目录下的文件夹列表
*/
func GetFolderList(c *gin.Context) {
	userInfo, _ := c.Get("userinfo")
	userId := userInfo.(jwt.MapClaims)["user_id"].(string)

	parentId := c.Param("folder_id")
	fmt.Println(parentId)
	if len(parentId) < 36 {
		if err := mysql.DB.Model(&models.Resource{}).Where("user_id = ? and name = ?", userId, userId).Select("id").Scan(&parentId).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "获取文件夹列表失败", "data": nil})
			return
		}
	}

	var folderList []schemas.FolderInfoResponse
	if err := mysql.DB.Model(&models.Resource{}).Where("user_id = ? and parent_id = ?", userId, parentId).Where("mime_type = ?", "folder").Select("id,name").Scan(&folderList).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "获取文件夹列表失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取文件夹列表成功", "data": folderList})
}
