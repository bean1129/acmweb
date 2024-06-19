package controller

import (
	"acmweb/app/authorization/model"
	"acmweb/app/authorization/service"
	"acmweb/app/errors"
	"acmweb/system"
	"acmweb/system/common"
	"fmt"
	"net/url"
	"reflect"

	"github.com/kataras/iris/v12"
	"github.com/mitchellh/mapstructure"
)

var Auth = new(AuthController)

type AuthController struct{}

// 授权申请控制器
// 一次可以绑定多个硬件指纹
// 申请时可以申请试用天数
func (c *AuthController) Apply(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	authReq := service.NewAuthService()
	if err := mapstructure.Decode(params, &authReq.ApplyData); err != nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	zs, err := authReq.Apply()
	if err != nil {
		zs.Code = errors.ErrOverApplyLimit.ErrorCode()
		zs.Msg = err.Error()
		// zs.Msg = ctx.Tr(errors.ErrOverApplyLimit.Error())
		system.Log.Error("Auth apply failed,errMsg:", zs.Msg)
	}
	ctx.JSON(zs)
}

func (c *AuthController) Approve(ctx iris.Context) {
	//逐级审批，审批通过后自动生成证书
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})

	authApprove := service.NewAuthService()
	if err := mapstructure.Decode(params, &authApprove.ApproveData); err != nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	zs, err := authApprove.Approve()
	if err != nil {
		system.Log.Error("Auth approve failed,errMsg:", zs.Msg)
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}

func (c *AuthController) ListAuth(ctx iris.Context) {
	zs := common.NewResult()
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	var a = model.AuthSearch{AuthState: -1, ApproveState: -1}
	if err := mapstructure.Decode(params, &a); err != nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	var authSearch = service.NewAuthService()
	zs, err := authSearch.List(a)
	if err != nil {
		system.Log.Error("ListApp  failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
func (c *AuthController) ListAuthNum(ctx iris.Context) {
	zs := common.NewResult()
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	var a = model.AuthSearch{AuthState: -1, ApproveState: -1}
	if err := mapstructure.Decode(params, &a); err != nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
	}
	var authSearch = service.NewAuthService()
	zs, err := authSearch.ListNum(a)
	if err != nil {
		system.Log.Error("ListApp  failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
func (c *AuthController) GetAuthInfoByApp(ctx iris.Context) {
	//查询指定应用的授权明细
}

func (c *AuthController) GetGrantedApp(ctx iris.Context) {
	//统计各个应用在各个组织的授权数量

}

func (c *AuthController) GetActiveApp(ctx iris.Context) {
	//统计各个单位正在使用的应用数量
}

func (c *AuthController) GetUsedCount(ctx iris.Context) {
	//统计各个应用在各个组织的使用次数
	zs := common.NewResult()
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	var appID = make([]int64, 0)
	for _, app := range params["app_id"].([]interface{}) {
		appID = append(appID, system.Text.Conv.Int64(app))
	}
	tierLevel := system.Text.Conv.Int(params["tier_level"])
	var authSearch = service.NewAuthService()
	zs, err := authSearch.AuthStatis(appID, tierLevel)
	if err != nil {
		system.Log.Error("GetUsedCount  failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
func (c *AuthController) GetStatisData(ctx iris.Context) {
	//按层级统计各个应用、按应用统计各个层级数据
	zs := common.NewResult()
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	var appID = make([]int64, 0)
	for _, app := range params["app_id"].([]interface{}) {
		appID = append(appID, system.Text.Conv.Int64(app))
	}
	tierLevel := system.Text.Conv.Int(params["tier_level"])
	var authSearch = service.NewAuthService()
	zs, err := authSearch.AuthStatisByAxis(appID, tierLevel)
	if err != nil {
		system.Log.Error("GetUsedCount  failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
func (c *AuthController) DownloadLicFile(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	authId := system.Text.Conv.Int64(params["auth_id"])
	var authSer = service.NewAuthService()
	srcFile, destFile, err := authSer.DownloadLicFile(authId)
	if err != nil {
		zs := common.NewResult()
		zs.Msg = ctx.Tr(err.Error())
		ctx.JSON(zs)
	} else {
		system.Log.Info("Download file, from ", srcFile, " to ", destFile)
		destFile = url.QueryEscape(destFile)
		ctx.Header("Access-Control-Expose-Headers", "Content-disposition")
		ctx.Header("Content-disposition", fmt.Sprintf("attachment;filename*=utf-8''%s", destFile))
		//ctx.Header("Content-disposition", fmt.Sprintf("attachment;filename=\"%s\"", destFile))
		ctx.SendFile(srcFile, destFile)
	}
}
func (c *AuthController) ListAuthMode(ctx iris.Context) {
	zs := common.NewResult()
	//zs.Data = append(zs.Data, model.AppType{
	//	1,
	//	"GV",
	//	"通用型",
	//})
	//zs.Data = append(zs.Data, model.AppType{
	//	2,
	//	"TV",
	//	"试用型",
	//})
	//zs.Data = append(zs.Data, model.AppType{
	//	3,
	//	"CV",
	//	"受控型",
	//})
	var authSearch = service.NewAuthService()
	zs, err := authSearch.ListAuthMode()
	if err != nil {
		system.Log.Error("ListApp  failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
func (auth *AuthController) DeleteAuth(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	ids := make([]int64, 1)
	v := reflect.ValueOf(params["auth_id"])
	if v.Kind() == reflect.Slice {
		idsparams := params["auth_id"].([]interface{})
		for _, v := range idsparams {
			ids = append(ids, int64(v.(float64)))
		}
	} else {
		ids[0] = int64(params["auth_id"].(float64))
	}
	var authSer = service.NewAuthService()
	zs, err := authSer.DeleteAuth(ids)
	if err != nil {
		system.Log.Error("Delete application failed,errMsg:", zs.Msg)
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
