package controller

import (
	"github.com/gookit/validate"
	"github.com/kataras/iris/v12"

	"acmweb/app/framework/service"
	"acmweb/app/framework/verification"
	"acmweb/constants"
	"acmweb/system/web"
)

var Login = new(LoginController)

type LoginController struct{}

func (c *LoginController) Login(ctx iris.Context) {
	if ctx.Method() != "POST" {
		return
	}
	data := ctx.Values().Get("data")
	if data == nil {
		return
	}
	params := data.(map[string]interface{})
	username := params["username"].(string)
	password := params["password"].(string)
	// 登录参数
	var req verification.LoginReq
	req.UserName = username
	req.Password = password
	// 参数校验
	v := validate.Struct(req)
	if !v.Validate() {
		ctx.JSON(web.JsonResult{
			Code: constants.ErrParamMiss,
			Msg:  v.Errors.One(),
		})
		return
	}
	// 系统登录
	token, err := service.Login.UserLogin(req.UserName, req.Password, ctx)
	if err != nil {
		// 登录错误
		ctx.JSON(web.JsonResult{
			Code: -1,
			Msg:  err.Error(),
		})
		return
	}
	tk := make([]map[string]interface{}, 0)
	tk = append(tk, iris.Map{
		"access_token": token,
	})
	// 登录成功
	ctx.JSON(web.JsonResult{
		Code: 0,
		Msg:  "登录成功",
		Data: tk,
	})
}

func (c *LoginController) Logout(ctx iris.Context) {
	if ctx.Method() != "POST" {
		return
	}
	_, err := service.Login.UserLogout(ctx)
	if err == nil {
		ctx.JSON(web.JsonResult{
			Code: 0,
			Msg:  "注销成功",
		})
		return
	} else {
		ctx.JSON(web.JsonResult{
			Code: -1,
			Msg:  "注销失败，原因：" + err.Error(),
		})
	}
}
