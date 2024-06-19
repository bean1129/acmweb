package router

import (
	"acmweb/app/organization/controller"
	"github.com/kataras/iris/v12"
)

func RegisterOrginazeRouter(app *iris.Application) {
	// 登录
	org := app.Party("/organization")
	{
		org.Post("/create", controller.ORG.Create)
		org.Post("/rename", controller.ORG.Rename)
		org.Post("/delete", controller.ORG.Delete)
		org.Post("/move", controller.ORG.Move)
		org.Post("/get", controller.ORG.List)
		org.Post("/getTier", controller.ORG.ListTier)
	}
}
