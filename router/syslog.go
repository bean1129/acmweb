package router

import (
	"github.com/kataras/iris/v12"

	"acmweb/app/framework/controller"
)

func RegisterSysLogRouter(app *iris.Application) {
	// 系统日志
	login := app.Party("/syslog")
	{
		login.Post("/create", controller.SysLog.CreateSysLog)
		login.Post("/list", controller.SysLog.ListLog)
	}
}
