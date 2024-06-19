package common

const (
	ErrProtoTimeExpire  = -1000 // 请求时间过期
	ErrProtoSignInvalid = -1001 // 请求签名无效
	ErrParamMiss        = -1002 // 参数丢失
	ErrDBExecute        = -1003 // 数据库操作错误
)
