package controller

import (
	"acmweb/app/errors"
	"acmweb/app/framework/model"
	"acmweb/app/framework/service"
	"acmweb/system"
	"acmweb/system/common"
	"github.com/kataras/iris/v12"
	"github.com/mitchellh/mapstructure"
)

var SysLog = new(SysLogController)

type SysLogController struct{}

func (auth *SysLogController) CreateSysLog(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	var a model.SysLog
	if err := mapstructure.Decode(params, &a); err != nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
	}
	zs, err := service.SysLog.CreateSysLog(a)
	if err != nil {
		system.Log.Error("Create application failed,errMsg:", zs.Msg)
		zs.Code = err.(*errors.AppErr).ErrorCode()
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}

func (c *SysLogController) ListLog(ctx iris.Context) {
	zs := common.NewResult()
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	var a model.SysLogSearch
	if err := mapstructure.Decode(params, &a); err != nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	zs, err := service.SysLog.List(a)
	if err != nil {
		system.Log.Error("List syslog  failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
