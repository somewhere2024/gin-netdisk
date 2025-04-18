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

/*
文件相关的接口
*/

// 单文件上传
func FileUpload(c *gin.Context) {
	/*
		- `file`: 文件内容（表单字段）
		- `parent_id`: "string" // 父目录ID（可选，默认根目录）
	*/
	file, _ := c.FormFile("file")
	parentId := c.DefaultPostForm("parent_id", "")

	if file == nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "文件为空", "data": nil})
		return
	}
	if parentId == "" {
		userInfo, exist := c.Get("userinfo")
		if !exist {
			c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "未授权用户", "data": nil})
			return
		}
		userId := userInfo.(jwt.MapClaims)["user_id"]
		path := "./data/" + userId.(string)
		if !utils.FolderExists(path) {
			_ = services.CreateFolder(userId.(string), userId.(string), "", path)
		}
		path = path + "/" + file.Filename
		err := services.CreateFile(userId.(string), file.Filename, parentId, path, file.Header.Get("Content-Type"), file.Size) //创建resources
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "创建文件资源失败", "data": nil})
			return
		}
		_ = c.SaveUploadedFile(file, path) //保存在本地
	} else {
		userInfo, exist := c.Get("userinfo")
		if !exist {
			c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "未授权用户", "data": nil})
			return
		}
		userId := userInfo.(jwt.MapClaims)["user_id"]

		path := services.GetParentPath(parentId) + "/" + file.Filename //获取父文件夹的路径

		rel := services.CreateFile(userId.(string), file.Filename, parentId, path, file.Header.Get("Content-Type"), file.Size)
		if rel != nil {
			c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "创建文件资源失败", "data": nil})
			return
		}
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "保存文件失败", "data": nil})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "上传成功", "data": nil})
}

// 获取目录的文件列表
func FileGet(c *gin.Context) {
	parentId := c.DefaultQuery("parent_id", "")
	userInfo, exist := c.Get("userinfo")
	if !exist {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "未授权用户", "data": nil})
		return
	}
	userId := userInfo.(jwt.MapClaims)["user_id"].(string)
	var fileInfo []schemas.FileInfoResponse
	if parentId == "" {
		if err := mysql.DB.Model(&models.Resource{}).Where("user_id = ? and name = ?", userId, userId).Select("id").Scan(&parentId).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "获取父目录失败", "data": nil})
			return
		}
	}
	if err := mysql.DB.Model(&models.Resource{}).Where("user_id = ? and parent_id = ?", userId, parentId).Select("id,name").Scan(&fileInfo).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "获取文件列表失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取文件列表成功", "data": fileInfo})
}

// 获取文件列表之后可以指定资源文件id进行download

/*
首先通过id从资源表中获取路径，顺便判断资源是否存在，通过路径将文件return给用户
*/
func FileDownload(c *gin.Context) {
	fileId := c.Param("file_id")
	if fileId == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "id为空", "data": nil})
		return
	}
	var downloadResponse schemas.FileDownloadResponse
	if err := mysql.DB.Model(&models.Resource{}).Where("id = ?", fileId).Select("path, name").Scan(&downloadResponse).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "获取文件路径失败", "data": nil})
		return
	}
	c.Header("Content-Type", "application/octet-stream")                           // 表示是文件流，唤起浏览器下载，一般设置了这个，就要设置文件名
	c.Header("Content-Disposition", "attachment; filename="+downloadResponse.Name) // 用来指定下载下来的文件名
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(downloadResponse.Path)
}

/*
文件删除思路：
通过路径参数获取文件id，
然后软删除，只更改资源的状态，不删除文件
并且把文件的的信息添加到回收站
*/
func FileDelete(c *gin.Context) {
	fileId := c.Param("file_id")
	userInfo, exists := c.Get("userinfo")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "未授权用户", "data": nil})
		return
	}
	var fileInfo models.Resource

	if err := mysql.DB.Where("id = ?", fileId).First(&fileInfo).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "文件不存在", "data": nil})
		return
	}
	userId := userInfo.(jwt.MapClaims)["user_id"].(string)
	if err := services.SoftDeleteFile(userId, fileInfo.ID, fileInfo.ParentId); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "文件删除失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "文件软删除成功", "data": nil})
}

/*
重命名文件
根据文件id，修改文件名

再对文件的路径进行修改
*/

func FileRename(c *gin.Context) {
	fileId := c.Param("file_id")
	newFileName := c.PostForm("new_name")
	file := models.Resource{}
	if err := mysql.DB.Where("id = ?", fileId).First(&file).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "文件不存在", "data": nil})
		return
	}
	file.Name = newFileName
	mysql.DB.Save(&file)

	//对文件进行重命名
	oldPath, _ := services.GetFilePath(fileId)
	newPath := services.GetParentPath(file.ParentId) + "/" + newFileName
	if err := services.RenameFile(oldPath, newPath); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "文件重命名失败", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "文件重命名成功", "data": nil})
}
