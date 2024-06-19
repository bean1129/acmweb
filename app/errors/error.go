package errors

import (
	appconst "acmweb/app/constants"
	"acmweb/constants"
	"acmweb/system/common"
)

type AppErr struct {
	code int
	err  string
}

func (e *AppErr) Error() string {
	return e.err
}
func (e *AppErr) ErrorCode() int {
	return e.code
}

var (
	ErrNo              = AppErr{common.OK, "Success"}
	ErrDBexecute       = AppErr{constants.ErrDBExecute, "DBExecuteFailed"}
	ErrParamMiss       = AppErr{constants.ErrParamMiss, "ParamMiss"}
	ErrNameExists      = AppErr{constants.ErrNameExists, "NameExists"}
	ErrDBQry           = AppErr{constants.ErrDBQry, "DBQueryFailed"}
	ErrParamType       = AppErr{constants.ErrParamType, "ParamTypeError"}
	ErrUserExists      = AppErr{appconst.ErrUserExists, "UserExists"}
	ErrNopeAppGRP      = AppErr{appconst.ErrNopeAppGRP, "NopeAppGrp"}
	ErrNopeAppWF       = AppErr{appconst.ErrNopeAppWF, "NopeAppWF"}
	ErrSameUserPwd     = AppErr{appconst.ErrSameUserPwd, "PasswordSame"}
	ErrInvaildApprover = AppErr{appconst.ErrInvaildApprover, "InvaildApprover"}
	ErrGenicense       = AppErr{appconst.ErrGenicense, "GenLicFailed"}
	ErrOverApplyLimit  = AppErr{appconst.ErrOverApplyLimit, "OverApplyLimit"}
)
