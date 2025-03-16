package xerr

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 错误全局定义
 * @Date: 2025-02-18 22:55
 */

// 定义一些业务常见的错误

var (
	UnDefinedError = NewAppError(00000, "undefind err")
	OtherError     = NewAppError(00001, "other err")

	UserNotExist = NewAppError(01001, "user is not exist")
	UserExists   = NewAppError(01002, "user already exists")
	TokenExpire  = NewAppError(01003, "token was expired")
	InviteError  = NewAppError(01004, "invite code is error")

	// Etcd 工具包
	EtcdBackGroundError   = NewAppError(02000, "etcd service err")            // Etcd 模块默认初始化错误
	EtcdErrNotInitialized = NewAppError(02001, "etcd client not initialized") // Etcd 未初始化
	EtcdErrInvalidConfig  = NewAppError(02002, "invalid etcd config ")        // Etcd 空配置
	// Redis 工具包
	RedisBackGroundError = NewAppError(0300, "etcd service err") // Redis默认初始错误
)
