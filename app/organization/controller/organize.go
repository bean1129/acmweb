package controller

import (
	"acmweb/app/errors"
	"acmweb/app/organization/model"
	"acmweb/app/organization/service"
	"acmweb/system"
	"acmweb/system/common"
	"acmweb/system/text"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

var ORG = new(OrganizeController)

type OrganizeController struct {
}

func (o *OrganizeController) Create(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}

	params := data.(map[string]interface{})
	u := &model.Unit{State: 1}
	if err := mapstructure.Decode(params, u); err != nil {
		system.Log.Error("Invaild params,err:", err.Error())
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	//创建
	zs, err := service.Org.Create(u)
	if err != nil {
		zs.Msg = ctx.Tr(err.Error())
		zs.Code = err.(*errors.AppErr).ErrorCode()
		system.Log.Error("Create unit failed,errMsg:", zs.Msg)
	}
	ctx.JSON(zs)
}

func (o *OrganizeController) Delete(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	ids := make([]int64, 1)
	//是否批量
	v := reflect.ValueOf(params["unit_id"])
	if v.Kind() == reflect.Slice {
		idsparams := params["unit_id"].([]interface{})
		for _, v := range idsparams {
			ids = append(ids, int64(v.(float64)))
		}
	} else {
		ids = append(ids, system.Text.Conv.Int64(params["user_id"]))
	}

	zs, err := service.Org.Delete(ids)
	if err != nil {
		system.Log.Error("Delete unit failed,errMsg:", zs.Msg)
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}

func (o *OrganizeController) Move(ctx iris.Context) {
	zs := common.NewResult()
	ctx.JSON(zs)
}

func (o *OrganizeController) Rename(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	m := model.Unit{UnitId: text.Conv.Int64(params["unit_id"]), UnitName: params["unit_name"].(string)}
	zs, err := service.Org.Rename(m)
	if err != nil {
		system.Log.Error("Rename unit failed,errMsg:", err.Error())
		zs.Code = err.(*errors.AppErr).ErrorCode()
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}

func (o *OrganizeController) List(ctx iris.Context) {
	zs := common.NewResult()
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	u := &model.UnitSearch{State: -1}
	if err := mapstructure.Decode(params, u); err != nil {
		fmt.Println(err)
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	zs, err := service.Org.List(u)
	if err != nil {
		system.Log.Error("getPageUsers  failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
func (o *OrganizeController) ListTier(ctx iris.Context) {
	zs := common.NewResult()
	zs, err := service.Org.ListTier()
	if err != nil {
		system.Log.Error("List tier failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
