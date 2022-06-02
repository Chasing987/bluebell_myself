package mysql

import "errors"

var (
	ErrorUserExit        = errors.New("用户已存在。。。。")
	ErrorUserNotExit     = errors.New("用户不存在！！！")
	ErrorInvalidPassword = errors.New("用户或密码错误")
	ErrorInvalidID       = errors.New("无效的id")
)
