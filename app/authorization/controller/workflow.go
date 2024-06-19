package controller

import (
	"acmweb/app/authorization/model"
	"acmweb/app/authorization/service"
	"acmweb/app/errors"
	"acmweb/constants"
	"acmweb/system"
	"acmweb/system/common"
	"reflect"

	"github.com/kataras/iris/v12"
	"github.com/mitchellh/mapstructure"
)

var AuthWF = new(AuthWFController)

type AuthWFController struct{}

func (wf *AuthWFController) CreateWorkFlow(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(constants.ErrParamMiss, constants.Message[constants.ErrParamMiss]))
		return
	}
	params := data.(map[string]interface{})
	var w model.WF
	if err := mapstructure.Decode(params, &w); err != nil {
		ctx.JSON(common.NewErrResult(constants.ErrParamMiss, constants.Message[constants.ErrParamMiss]))
		return
	}

	zs, err := service.AuthWF.CreateWorkFlow(w)
	if err != nil {
		zs.Msg = ctx.Tr(err.Error())
		system.Log.Error("Create application failed,errMsg:", zs.Msg)
	}
	ctx.JSON(zs)
}
func (wf *AuthWFController) GetWorkFlow(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(constants.ErrParamMiss, constants.Message[constants.ErrParamMiss]))
		return
	}
	params := data.(map[string]interface{})
	grpId := system.Text.Conv.Int64(params["group_id"])
	start := 1
	pageSize := 10
	if val, ok := params["page_num"]; ok {
		start = system.Text.Conv.Int(val)
	}
	if val, ok := params["page_size"]; ok {
		pageSize = system.Text.Conv.Int(val)
	}
	start = (start - 1) * pageSize
	zs, err := service.AuthWF.GetWorkFlow(start, pageSize, grpId)
	if err != nil {
		zs.Msg = ctx.Tr(err.Error())
		system.Log.Error("Get work flow information failed,errMsg:", zs.Msg)
	}
	ctx.JSON(zs)
}
func (auth *AuthWFController) Delete(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	ids := make([]int64, 1)
	v := reflect.ValueOf(params["wf_id"])
	if v.Kind() == reflect.Slice {
		idsparams := params["wf_id"].([]interface{})
		for _, v := range idsparams {
			ids = append(ids, int64(v.(float64)))
		}
	} else {
		ids[0] = int64(params["wf_id"].(float64))
	}

	zs, err := service.AuthWF.Delete(ids)
	if err != nil {
		system.Log.Error("Delete application failed,errMsg:", zs.Msg)
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
func (auth *AuthWFController) Rename(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	var a model.WF
	if err := mapstructure.Decode(params, &a); err != nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
	}
	zs, err := service.AuthWF.RenameWorkFlow(a)
	if err != nil {
		system.Log.Error("Rename application failed,errMsg:", zs.Msg)
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
func (wf *AuthWFController) GetWFLog(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(constants.ErrParamMiss, constants.Message[constants.ErrParamMiss]))
		return
	}
	params := data.(map[string]interface{})
	wfId := system.Text.Conv.Int64(params["wf_inst_id"])

	zs, err := service.AuthWF.GetWFLog(wfId)
	if err != nil {
		zs.Msg = ctx.Tr(err.Error())
		system.Log.Error("Get work flow information failed,errMsg:", zs.Msg)
	}
	ctx.JSON(zs)
}
func (wf *AuthWFController) GetHisAppv(ctx iris.Context) {
	zs := common.NewResult()
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	var a = model.WFHisSearch{AppvState: -1, ApproverId: system.Text.Conv.Int64(params["approver_id"])}
	if err := mapstructure.Decode(params, &a); err != nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
	}
	zs, err := service.AuthWF.GetHisAppv(a)
	if err != nil {
		system.Log.Error("ListApp  failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
