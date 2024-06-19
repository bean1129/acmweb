package service

import (
	"errors"
	"github.com/kataras/iris/v12"

	"acmweb/system"
)

var (
	Login = new(loginService)
)

type loginService struct{}

// UserLogin 系统登录
func (s *loginService) UserLogin(username, password string, ctx iris.Context) (string, error) {
	// 查询用户
	ds, err := system.Data.DB.Query("select user_id,user_name,passwd,state from user where user_code=?", username)
	if ds == nil || err != nil {
		return "", errors.New("用户名或者密码不正确")
	}
	// 密码校验
	pwd, _ := system.Secret.MD5.Password(password + system.Config.Application.SecretKey + username)
	if ds["passwd"] != pwd {
		return "", errors.New("密码不正确")
	}
	// 判断当前用户状态
	if system.Text.Conv.Int(ds["state"]) != 1 {
		return "", errors.New("您的账号已被禁用,请联系管理员")
	}
	// 更新登录时间、登录IP
	system.Data.DB.Query("update user set lasted_login_time=now(),lasted_ip=? where user_code=?", ctx.Host(), username)
	// 生成Token
	token, _ := system.Web.Token.Generate(system.Text.Conv.Int(ds["user_id"]), username, pwd)
	session := system.Data.Session.Start(ctx)
	session.Set("authenticated", true)
	session.Set("userid", username)

	system.Log.Infof("[%v] Token: [%s]", ctx.Values().Get("tranId"), token)

	// 返回token
	return token, nil
}

// UserLogout 用户退出登录状态
func (s *loginService) UserLogout(ctx iris.Context) (string, error) {
	session := system.Data.Session.Start(ctx)

	// 移除一个认证用户
	session.Set("authenticated", false)
	session.Set("userid", "")
	// 或者摧毁 session:
	session.Destroy()

	return "", nil
}

func (s *loginService) GetProfile(ctx iris.Context) (string, error) {
	return "", errors.New("配置文件")
}
