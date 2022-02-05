package v1

import (
	"github.com/DMwangnima/easy-disk/metadata/user"
	"github.com/gin-gonic/gin"
)

var (
	UserManager user.Manager
)

func Login(ctx *gin.Context) {

}

func Logout(ctx *gin.Context) {

}

// /user/:uid
func GetUser(ctx *gin.Context) {
    id := ctx.Param("id")
    if id == "" {
    	// return 400
		return
	}
	resUser, err := UserManager.Get(ctx, id)
	if err != nil {
		// 根据返回错误类型，返回相应状态码
		return
	}
    ctx.JSON(200, resUser)
}

func CreateUser(ctx *gin.Context) {

}

func UpdateUser(ctx *gin.Context) {

}

// DeleteUser
// Not exposed to normal user
func DeleteUser(ctx *gin.Context) {

}
