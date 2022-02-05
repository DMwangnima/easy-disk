package user

// UserOption 在原有的FunctionalOptions上做了改进，方便拿到字段名
type UserOption func(*User) string

func WithUserNickName(name string) UserOption {
	return func(user *User) string {
		user.NickName = name
		return "nick_name"
	}
}

func WithUserEmail(email string) UserOption {
	return func(user *User) string {
		user.Email = email
		return "email"
	}
}

//func WithPassword(password string) UserOption {
//	return func(user *User) {
//		user.RawPassword = password
//	}
//}
