package constants

const (
	ErrProtoTimeExpire  = -1000 // 请求时间过期
	ErrProtoSignInvalid = -1001 // 请求签名无效
	ErrParamMiss        = -1002 // 参数丢失
	ErrDBExecute        = -1003 // 数据库操作错误
	ErrParamType        = -1004
	ErrDBQry            = -1005
)

const (
	DBTIMEFORMAT = "2006-01-02 15:04:05"
	DBDATEFORMAT = "2006-01-02"
)
