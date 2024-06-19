package constants

import "acmweb/system/common"

var (
	Message = make(map[int]string)
)

func init() {
	Message[common.OK] = "Success"
	Message[ErrDBExecute] = "DBExecuteFailed"
	Message[ErrParamMiss] = "ParamMiss"
	Message[ErrNameExists] = "NameExists"
}
