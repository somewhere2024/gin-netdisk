package services

import (
	"gin-netdisk/internal/dao/mysql"
	"gin-netdisk/internal/models"
	"gin-netdisk/internal/utils"
	"os"
)

func UserCreateFolder(user_id string) error {
	//创建数据库中的文件夹资源
	resource := models.Resource{
		UserId:   user_id,
		Name:     user_id,
		Type:     models.Folder,
		ParentId: "",
		MimeType: "folder",
		Path:     "./data/" + user_id,
	}
	mysql.DB.Create(&resource)
	//创建文件夹
	path := "./data/" + user_id
	err := os.MkdirAll(path, os.ModePerm) //文件夹权限为777
	if err != nil {
		utils.Logger.Warn("用户的文件夹创建失败")
		return err
	}
	return nil
}

func CreateFolder(userId, name, parentId, path string) error {
	//创建文件夹资源
	resource := models.Resource{
		UserId:   userId,
		Name:     name,
		Type:     models.Folder,
		ParentId: parentId,
		MimeType: "folder",
		Path:     path,
	}
	err := mysql.DB.Create(&resource).Error
	if err != nil {
		utils.Logger.Warn("创建文件夹资源失败")
		return err
	}
	// 创建文件夹
	err = os.MkdirAll(path, os.ModePerm) //文件夹权限为777
	if err != nil {
		utils.Logger.Warn("用户的文件夹创建失败")
		return err
	}
	return nil
}
