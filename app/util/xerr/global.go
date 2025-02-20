package xerr

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 错误全局定义
 * @Date: 2025-02-18 22:55
 */

// 定义一些业务常见的错误
var (
	UnDefinedError = NewAppError(0000, "undefind err")
	UserNotExist   = NewAppError(1001, "user is not exist")
	UserExists     = NewAppError(1002, "user already exists")
	InviteError    = NewAppError(1003, "invite code is error")
)
