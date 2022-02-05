package v1

import (
	v1 "github.com/DMwangnima/easy-disk/data/api/v1"
	"github.com/gin-gonic/gin"
)

const (
	STORAGE_PATH = "storage"
)

func InitStorage(router *gin.RouterGroup) {
	storageGroup := router.Group(STORAGE_PATH)
	{
		storageGroup.GET("", v1.Get)
		storageGroup.POST("", v1.Put)
	}
}
