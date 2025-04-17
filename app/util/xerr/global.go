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
	OtherError     = NewAppError(0001, "other err")

	UserNotExist = NewAppError(1001, "user is not exist")
	UserExists   = NewAppError(1002, "user already exists")
	TokenExpire  = NewAppError(1003, "token was expired")
	InviteError  = NewAppError(1004, "invite code is error")

	// Etcd 工具包
	EtcdBackGroundError         = NewAppError(2000, "etcd service err")            // Etcd 模块默认初始化错误
	EtcdErrNotInitialized       = NewAppError(2001, "etcd client not initialized") // Etcd 未初始化
	EtcdErrInvalidConfig        = NewAppError(2002, "invalid etcd config ")        // Etcd 空配置
	EtcdErrNoKey                = NewAppError(2003, "no key provided")
	EtcdErrNoValue              = NewAppError(2004, "no value provided")
	EtcdErrUnknownMethod        = NewAppError(2005, "unknown method")
	EtcdErrConfigSplitBrain     = NewAppError(2006, "nodes will be split brain") // Etcd 偶数节点容易造成脑裂
	EtcdUnConfTLSOption         = NewAppError(2007, "undefined tls option")
	EtcdErrClientNotInitialized = NewAppError(2008, "etcd client is not inited")
	// Redis 工具包
	RedisBackGroundError = NewAppError(3000, "etcd service err") // Redis默认初始错误

)
