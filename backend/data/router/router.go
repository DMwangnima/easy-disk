package router

import (
	v1 "github.com/DMwangnima/easy-disk/data/router/v1"
	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	router := gin.Default()
	v1Group := router.Group("v1")
	{
		v1.InitStorage(v1Group)
	}
	return router
}