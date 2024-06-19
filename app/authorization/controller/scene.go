package controller

import (
	"acmweb/app/authorization/model"
	"acmweb/app/authorization/service"
	"acmweb/app/errors"
	"acmweb/system"
	"acmweb/system/common"
	"github.com/kataras/iris/v12"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

var Scene = new(SceneController)

type SceneController struct{}

func (auth *SceneController) CreateScene(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	var a model.Scene

	if err := mapstructure.Decode(params, &a); err != nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
	}
	zs, err := service.Scene.CreateScene(a)
	if err != nil {
		system.Log.Error("Create scene failed,errMsg:", zs.Msg)
		zs.Code = err.(*errors.AppErr).ErrorCode()
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}

func (auth *SceneController) DeleteScene(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	ids := make([]int64, 1)
	v := reflect.ValueOf(params["inst_id"])
	if v.Kind() == reflect.Slice {
		idsparams := params["inst_id"].([]interface{})
		for _, v := range idsparams {
			ids = append(ids, int64(v.(float64)))
		}
	} else {
		ids[0] = int64(params["inst_id"].(float64))
	}

	zs, err := service.Scene.DeleteScene(ids)
	if err != nil {
		system.Log.Error("Delete scene failed,errMsg:", zs.Msg)
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
func (auth *SceneController) CreateGroup(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	var a model.SceneGroup

	if err := mapstructure.Decode(params, &a); err != nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
	}
	zs, err := service.Scene.CreateGrp(a)
	if err != nil {
		system.Log.Error("Create scene group failed,errMsg:", zs.Msg)
		zs.Code = err.(*errors.AppErr).ErrorCode()
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
func (auth *SceneController) DeleteGroup(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	ids := make([]int64, 1)
	v := reflect.ValueOf(params["grp_id"])
	if v.Kind() == reflect.Slice {
		idsparams := params["grp_id"].([]interface{})
		for _, v := range idsparams {
			ids = append(ids, int64(v.(float64)))
		}
	} else {
		ids[0] = int64(params["grp_id"].(float64))
	}

	zs, err := service.Scene.DeleteGroup(ids)
	if err != nil {
		system.Log.Error("Delete scene group failed,errMsg:", zs.Msg)
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
func (auth *SceneController) ListScene(ctx iris.Context) {
	zs := common.NewResult()
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	a := model.NewScenseSearch()
	if err := mapstructure.Decode(params, &a); err != nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
	}

	zs, err := service.Scene.ListScene(a)
	if err != nil {
		system.Log.Error("ListApp  failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
func (auth *SceneController) ListGroup(ctx iris.Context) {
	zs := common.NewResult()
	zs, err := service.Scene.ListGroup()
	if err != nil {
		system.Log.Error("List app group  failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
func (auth *SceneController) ListSceneApp(ctx iris.Context) {
	zs := common.NewResult()
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	ids := make([]int64, 1)
	v := reflect.ValueOf(params["inst_id"])
	if v.Kind() == reflect.Slice {
		idsparams := params["inst_id"].([]interface{})
		for _, v := range idsparams {
			ids = append(ids, int64(v.(float64)))
		}
	} else {
		ids[0] = int64(params["inst_id"].(float64))
	}

	zs, err := service.Scene.ListSceneApp(ids)
	if err != nil {
		system.Log.Error("ListApp  failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
