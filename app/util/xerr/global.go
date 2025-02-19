package xerr

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 错误全局定义
 * @Date: 2025-02-18 22:55
 */

// 定义一些业务常见的错误
var (
	UnDefinedErr = &AppError{Code: 0000, Message: "undefind err"}
	UserNotExist = &AppError{Code: 1001, Message: "user is not exist"}
	UserExists   = &AppError{Code: 1002, Message: "user already exists"}
)
