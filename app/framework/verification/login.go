package verification

import "github.com/gookit/validate"

// LoginReq 系统登录
type LoginReq struct {
	UserName string `form:"username" verification:"required"`
	Password string `form:"password" verification:"required"`
}

// Messages 登录参数校验
func (v LoginReq) Messages() map[string]string {
	return validate.MS{
		"UserName.required": "登录用户名不能为空.",
		"Password.required": "登录密码不能为空.",
	}
}
