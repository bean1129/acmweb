package router

import (
	"github.com/kataras/iris/v12"

	"acmweb/app/framework/controller"
)

func RegisterLoginRouter(app *iris.Application) {
	// 登录
	login := app.Party("/")
	{
		login.Get("/", controller.Login.Login)
		login.Post("/login", controller.Login.Login)
		login.Post("/logout", controller.Login.Logout)
	}
}
