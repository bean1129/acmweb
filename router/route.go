package router

import (
	"github.com/kataras/iris/v12"
)

func RegisterRouter(app *iris.Application) {
	RegisterLoginRouter(app)
	RegisterPermissionRouter(app)
	RegisterAuthRouter(app)
	RegisterOrginazeRouter(app)
	RegisterSysLogRouter(app)
}
