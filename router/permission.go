package router

import (
	"github.com/kataras/iris/v12"

	"acmweb/app/permission/controller"
)

func RegisterPermissionRouter(app *iris.Application) {
	// 登录
	permission := app.Party("/permission")
	{
		//角色
		role := permission.Party("/roles")
		role.Post("/deleteRole", controller.Role.Delete)
		role.Post("/createRole", controller.Role.Create)
		role.Post("/getAllRoles", controller.Role.List)
		role.Post("/getRolesByCond", controller.Role.Find)
		role.Post("/renameRole", controller.Role.Rename)
		role.Post("/setPermission", controller.Role.SetPermission)
		role.Post("/setRoleFunc", controller.Role.SetRoleFunc)
		role.Post("/getRoleFunc", controller.Role.GetRoleFunc)
		role.Post("/setRoleState", controller.Role.UpdateState)
		role.Post("/getAllFunc", controller.Role.ListFunc)
	}
	{
		//用户
		user := permission.Party("/users")
		user.Post("/getUserInfo", controller.User.GetUsersInfo)
		user.Post("/getPageUsers", controller.User.GetUsersInfo)
		user.Post("/createUser", controller.User.Create)
		user.Post("/deleteUser", controller.User.Delete)
		user.Post("/renameUser", controller.User.RenameUser)
		user.Post("/setpasswd", controller.User.SetPasswd)
		user.Post("/setUserState", controller.User.UpdateState)
		user.Post("/modify", controller.User.ModifyUser)
		user.Post("/setUserRole", controller.User.SetUserRole)
	}
}
