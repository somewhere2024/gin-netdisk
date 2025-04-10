package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func TestStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": "200", "message": "ok"})
}
