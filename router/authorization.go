package router

import (
	"github.com/kataras/iris/v12"

	"acmweb/app/authorization/controller"
)

func RegisterAuthRouter(app *iris.Application) {
	//权限管理
	auth := app.Party("/auth")
	{
		auth.Post("/apply", controller.Auth.Apply)
		auth.Post("/approve", controller.Auth.Approve)
		auth.Post("/getAuthReqs", controller.Auth.ListAuth)
		auth.Post("/getAuthReqsNum", controller.Auth.ListAuthNum)
		auth.Post("/getAuthInfo", controller.Auth.GetAuthInfoByApp)
		auth.Post("/getGrantedApp", controller.Auth.GetGrantedApp)
		auth.Post("/getActiveApp", controller.Auth.GetActiveApp)
		auth.Post("/getUsedCount", controller.Auth.GetUsedCount)
		auth.Post("/download", controller.Auth.DownloadLicFile)
		auth.Post("/authMode", controller.Auth.ListAuthMode)
		auth.Post("/delete", controller.Auth.DeleteAuth)
	}
	{
		//应用
		_app := auth.Party("/app")
		_app.Post("/create", controller.AuthApp.CreateApp)
		_app.Post("/rename", controller.AuthApp.RenameApp)
		_app.Post("/delete", controller.AuthApp.DeleteApp)
		_app.Post("/modify", controller.AuthApp.ModifyApp)
		_app.Post("/setModule", controller.AuthApp.ChangeModule)
		_app.Post("/getModules", controller.AuthApp.ListModule)
		_app.Post("/deleteModules", controller.AuthApp.DeleteModule)
		_app.Post("/getApps", controller.AuthApp.ListApp)
		_app.Post("/getAppType", controller.AuthApp.ListAppType)
		_app.Post("/getAppGrp", controller.AuthApp.ListAppGrp)
		_app.Post("/download", controller.AuthApp.DownloadApp)
	}
	{
		//流程
		_wf := auth.Party("/wf")
		_wf.Post("/create", controller.AuthWF.CreateWorkFlow)
		_wf.Post("/rename", controller.AuthWF.Rename)
		_wf.Post("/get", controller.AuthWF.GetWorkFlow)
		_wf.Post("/delete", controller.AuthWF.Delete)
		_wf.Post("/getWFLog", controller.AuthWF.GetWFLog)
		_wf.Post("/getHisAppv", controller.AuthWF.GetHisAppv)
	}
	{
		//场景
		_scene := auth.Party("/scene")
		_scene.Post("/create", controller.Scene.CreateScene)
		_scene.Post("/delete", controller.Scene.DeleteScene)
		_scene.Post("/createGroup", controller.Scene.CreateGroup)
		_scene.Post("/deleteGroup", controller.Scene.DeleteGroup)
		_scene.Post("/list", controller.Scene.ListScene)
		_scene.Post("/listGroup", controller.Scene.ListGroup)
		_scene.Post("/listApp", controller.Scene.ListSceneApp)
	}
}
