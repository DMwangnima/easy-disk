package v1

import (
	"github.com/DMwangnima/easy-disk/metadata/api/v1"
	"github.com/gin-gonic/gin"
)

const (
	USER_PATH = "users"
)

func InitUser(router *gin.RouterGroup) {
	userGroup := router.Group(USER_PATH)
	{
		userGroup.POST("login", v1.Login)
		userGroup.POST("logout", v1.Logout)
		userGroup.GET(":id", v1.GetUser)
		userGroup.POST("", v1.CreateUser)
		userGroup.PATCH(":id", v1.UpdateUser)
		userGroup.DELETE(":id", v1.DeleteUser)
	}
}
