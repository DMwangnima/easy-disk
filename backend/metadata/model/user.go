package model

import (
	"github.com/DMwangnima/easy-disk/metadata/user"
	"time"
)

// User 业务实体
type User struct {
	Uid               uint64 `ddb:"uid"`
	NickName          string `ddb:"nick_name"`
	Email             string `ddb:"email"`
	EncryptedPassword string `ddb:"encrypted_password"`
	Group             uint8 `ddb:"group"`
	FileNum           uint64 `ddb:"file_num"`
	RegisterDate      time.Time `ddb:"register_date"`
	LastLoginDate     time.Time `ddb:"last_login_date"`
}

type GetUserResp struct {
	*user.User
}

type CreateUserReq struct {

}

type CreateUserResp struct {

}
