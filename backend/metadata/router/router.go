package router

import (
	v12 "github.com/DMwangnima/easy-disk/metadata/api/v1"
	v1 "github.com/DMwangnima/easy-disk/metadata/router/v1"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	engine := gin.Default()
    apiGroup := engine.Group("api")
    apiGroup.POST("login", v12.Login)
    apiGroup.POST("logout", v12.Logout)
    v1Group := apiGroup.Group("v1")
    {
    	v1.InitUser(v1Group)
	}
    return engine
}
