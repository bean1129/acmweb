package controller

import (
	"acmweb/app/authorization/model"
	"acmweb/app/authorization/service"
	"acmweb/app/errors"
	"acmweb/system"
	"acmweb/system/common"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/mitchellh/mapstructure"
	"io"
	"net/url"
	"os"
	"reflect"
)

var AuthApp = new(AuthAppController)

type AuthAppController struct{}

func (auth *AuthAppController) CreateApp(ctx iris.Context) {
	appFile, appFileHeader, _ := ctx.FormFile("app-file")
	data := ctx.Values().Get("data")
	params := data.(map[string]interface{})

	var a model.App
	if err := mapstructure.Decode(params, &a); err != nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
	}
	if appFileHeader != nil {
		addPath := "/app/"
		a.FileName = addPath + appFileHeader.Filename
		fullFilePath := system.Config.Attachment.FilePath + addPath
		if !system.Text.File.Exists(fullFilePath) {
			system.Text.File.Mkdir(fullFilePath)
		}
		fullFileName := fullFilePath + appFileHeader.Filename
		out, err := os.OpenFile(fullFileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
			return
		}
		defer out.Close()
		io.Copy(out, appFile)
	}

	zs, err := service.AuthApp.CreateApp(a)
	if err != nil {
		system.Log.Error("Create application failed,errMsg:", zs.Msg)
		zs.Code = err.(*errors.AppErr).ErrorCode()
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}

func (auth *AuthAppController) RenameApp(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	var a model.App
	if err := mapstructure.Decode(params, &a); err != nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
	}
	zs, err := service.AuthApp.RenameApp(a)
	if err != nil {
		system.Log.Error("Rename application failed,errMsg:", zs.Msg)
		zs.Code = err.(*errors.AppErr).ErrorCode()
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}

func (auth *AuthAppController) DeleteApp(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	ids := make([]int64, 1)
	v := reflect.ValueOf(params["app_id"])
	if v.Kind() == reflect.Slice {
		idsparams := params["app_id"].([]interface{})
		for _, v := range idsparams {
			ids = append(ids, int64(v.(float64)))
		}
	} else {
		ids[0] = int64(params["app_id"].(float64))
	}

	zs, err := service.AuthApp.DeleteApp(ids)
	if err != nil {
		system.Log.Error("Delete application failed,errMsg:", zs.Msg)
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}

func (auth *AuthAppController) ModifyApp(ctx iris.Context) {
	appFile, appFileHeader, _ := ctx.FormFile("app-file")
	data := ctx.Values().Get("data")
	params := data.(map[string]interface{})

	var a model.App
	if err := mapstructure.Decode(params, &a); err != nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
	}

	if appFileHeader != nil {
		addPath := "/app/"
		a.FileName = addPath + appFileHeader.Filename
		fullFilePath := system.Config.Attachment.FilePath + addPath
		if !system.Text.File.Exists(fullFilePath) {
			system.Text.File.Mkdir(fullFilePath)
		}
		fullFileName := fullFilePath + appFileHeader.Filename
		out, err := os.OpenFile(fullFileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
			return
		}
		defer out.Close()
		io.Copy(out, appFile)
	}

	zs, err := service.AuthApp.ModifyApp(a)
	if err != nil {
		system.Log.Error("Update application failed,errMsg:", zs.Msg)
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}

func (auth *AuthAppController) ChangeModule(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	var a model.AppModule
	if err := mapstructure.Decode(params, &a); err != nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	zs, err := service.AuthApp.ChangeModule(a)
	if err != nil {
		system.Log.Error("Create application failed,errMsg:", zs.Msg)
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}

func (auth *AuthAppController) ListApp(ctx iris.Context) {
	zs := common.NewResult()
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	var a model.AppSearch
	if err := mapstructure.Decode(params, &a); err != nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
	}

	zs, err := service.AuthApp.List(a)
	if err != nil {
		system.Log.Error("ListApp  failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}

func (auth *AuthAppController) ListAppType(ctx iris.Context) {
	zs := common.NewResult()
	zs, err := service.AuthApp.ListAppType()
	if err != nil {
		system.Log.Error("ListApp  failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
func (auth *AuthAppController) ListAppGrp(ctx iris.Context) {
	zs := common.NewResult()
	zs, err := service.AuthApp.ListAppGrp()
	if err != nil {
		system.Log.Error("List app group  failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}

func (auth *AuthAppController) ListModule(ctx iris.Context) {
	zs := common.NewResult()
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	appId := system.Text.Conv.Int64(params["app_id"])
	zs, err := service.AuthApp.ListModules(appId)
	if err != nil {
		system.Log.Error("ListApp  failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}

func (auth *AuthAppController) DeleteModule(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	ids := make([]int64, 1)
	v := reflect.ValueOf(params["id"])
	if v.Kind() == reflect.Slice {
		idsparams := params["id"].([]interface{})
		for _, v := range idsparams {
			ids = append(ids, int64(v.(float64)))
		}
	} else {
		ids[0] = int64(params["id"].(float64))
	}

	zs, err := service.AuthApp.DeleteModule(ids)
	if err != nil {
		system.Log.Error("Delete application failed,errMsg:", zs.Msg)
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}

func (auth *AuthAppController) DownloadApp(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	appId := system.Text.Conv.Int64(params["app_id"])
	srcFile, destFile, err := service.AuthApp.DownloadApp(appId)
	if err != nil {
		zs := common.NewResult()
		zs.Msg = ctx.Tr(err.Error())
		ctx.JSON(zs)
	} else {
		system.Log.Info("Download file :", destFile)
		destFile = url.QueryEscape(destFile)
		ctx.Header("Access-Control-Expose-Headers", "Content-disposition")
		ctx.Header("Content-disposition", fmt.Sprintf("attachment;filename*=utf-8''%s", destFile))
		//ctx.Header("Content-disposition", fmt.Sprintf("attachment;filename=\"%s\"", destFile))
		ctx.SendFile(srcFile, destFile)
	}
}
