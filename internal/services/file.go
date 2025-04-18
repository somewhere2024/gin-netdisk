package services

import (
	"gin-netdisk/internal/dao/mysql"
	"gin-netdisk/internal/models"
	"gin-netdisk/internal/utils"
	"os"
)

func CreateFile(userId, name, parentId, path, mimeType string, size int64) error {
	if parentId == "" {
		x := models.Resource{}
		_ = mysql.DB.Where("name = ? and user_id = ?", userId, userId).First(&x)
		parentId = x.ID
	}

	file := models.Resource{
		UserId:   userId,
		Name:     name,
		Type:     models.File,
		ParentId: parentId,
		MimeType: mimeType,
		Path:     path,
		Size:     size,
	}
	if err := mysql.DB.Create(&file).Error; err != nil {
		utils.Logger.Warn("创建文件资源失败")
		return err
	}

	return nil
}

func GetParentPath(parentId string) string {
	if parentId == "" {
		return ""
	}
	x := models.Resource{}
	err := mysql.DB.Where("id = ?", parentId).First(&x).Error
	if err != nil {
		utils.Logger.Warn("获取父级文件夹的路径失败")
		return ""
	}
	return x.Path
}

// 首先更改资源为删除状态
// 然后创建回收站资源
func SoftDeleteFile(userId, fileId, ParentId string) error {

	fileResource := models.Resource{}

	mysql.DB.Where("id = ?", fileId).First(&fileResource)
	fileResource.IsDelete = true
	mysql.DB.Save(&fileResource)
	trash := models.Trash{
		ResourceID:       fileId,
		UserID:           userId,
		OriginalParentId: ParentId,
	}
	if err := mysql.DB.Create(&trash).Error; err != nil {
		utils.Logger.Warn("创建回收站资源失败")
		return err
	}
	return nil
}

// 根据文件id返回文件路径
func GetFilePath(fileId string) (string, error) {
	var path string
	if err := mysql.DB.Model(&models.Resource{}).Where("id = ?", fileId).Select("path").Scan(&path).Error; err != nil {
		utils.Logger.Warn("获取文件路径失败")
		return "", err
	}
	return path, nil
}

func RenameFile(oldPathName, newPathName string) error {
	if oldPathName == newPathName {
		return nil
	}

	if err := os.Rename(oldPathName, newPathName); err != nil {
		return err
	}
	return nil
}
