package user

import (
	"context"
	"time"
)

const (
	NORMAL = 2 << iota
	ADMINISTRATOR
)

type User struct {
	Id                uint64    `ddb:"id"`
	NickName          string    `ddb:"nick_name"`
	Email             string    `ddb:"email"`
	EncryptedPassword string    `ddb:"encrypted_password"`
	Group             uint8     `ddb:"group"`
	FileNum           uint64    `ddb:"file_num"`
	RegisterDate      time.Time `ddb:"register_date"`
	LastLoginDate     time.Time `ddb:"last_login_date"`
}

type Manager interface {
	Get(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, opts ...UserOption) error
	Update(ctx context.Context, id string, opts ...UserOption) error
	Delete(ctx context.Context, id string) error
}
