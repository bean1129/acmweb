package verification

import (
	"github.com/gookit/validate"
	"regexp"
)

// 用户信息生成校验
type UserReq struct {
	UserName string `validate:"required"`
	UserCode string `validate:"required"`
	Password string `validate:"required"`
	Phone    string `validate:"PhoneValidator"`
}

// 参数校验
func (v UserReq) Messages() map[string]string {
	return validate.MS{
		"UserName.required":    "用户名称不能为空.",
		"UserCode.required":    "用户名不能为空.",
		"Password.required":    "登录密码不能为空.",
		"Phone.PhoneValidator": "手机号码不存在",
	}
}

func (v UserReq) PhoneValidator(val string) bool {
	bPass, _ := regexp.MatchString("^(13[0-9]|14[5|7]|15[0|1|2|3|5|6|7|8|9]|18[0|1|2|3|5|6|7|8|9])\\d{8}$", val)
	return bPass
}

type IpReq struct {
	IpAddress string `validate:"IpValidator"`
}

func (v IpReq) Messages() map[string]string {
	return validate.MS{
		"IpAddress.IpValidator": "无效的ip和端口格式",
	}
}

func (v IpReq) IpValidator(val string) bool {
	bPass, _ := regexp.MatchString("\\d+\\.\\d+\\.\\d+\\.\\d+:\\d+", val)
	return bPass
}

//手机号
type PhoneReq struct {
	Phone string `validate:"PhoneValidator"`
}

func (v PhoneReq) Messages() map[string]string {
	return validate.MS{
		"Phone.PhoneValidator": "手机号码不存在",
	}
}

func (v PhoneReq) PhoneValidator(val string) bool {
	bPass, _ := regexp.MatchString("^(13[0-9]|14[5|7]|15[0|1|2|3|5|6|7|8|9]|18[0|1|2|3|5|6|7|8|9])\\d{8}$", val)
	return bPass
}
