package controller

import (
	"acmweb/app/errors"
	"github.com/kataras/iris/v12"
	"github.com/mitchellh/mapstructure"
	"reflect"

	"acmweb/app/permission/model"
	"acmweb/app/permission/service"
	"acmweb/constants"
	"acmweb/system"
	"acmweb/system/common"
	"acmweb/system/text"
)

var Role = new(RoleController)

type RoleController struct{}

func (c *RoleController) Create(ctx iris.Context) {
	zs := common.NewResult()
	data := ctx.Values().Get("data")
	if data == nil {
		zs.Code = constants.ErrParamMiss
		zs.Msg = ctx.Tr(constants.Message[constants.ErrParamMiss])
		ctx.JSON(zs)
		return
	}
	params := data.(map[string]interface{})
	val := reflect.ValueOf(params["remark"])
	remark := ""
	if val.Kind() == reflect.String {
		remark = val.String()
	}
	m := model.Role{
		Id:     system.Common.UUID.NextVal(),
		Name:   params["name"].(string),
		PId:    text.Conv.Int64(params["pid"]),
		State:  1,
		Remark: remark,
		MId:    text.Conv.Int64(params["mid"]),
	}
	zs = service.Role.CreateRole(m)
	if zs.Code != 0 {
		system.Log.Error("Create role failed,errMsg:", zs.Msg)
	}
	zs.Msg = ctx.Tr(constants.Message[zs.Code])
	ctx.JSON(zs)
}

func (c *RoleController) Delete(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		return
	}
	params := data.(map[string]interface{})
	id := text.Conv.Int64(params["id"])
	zs := service.Role.DeleteRole(id)
	ctx.JSON(zs)
}

func (c *RoleController) Update(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	roleId := text.Conv.Int64(params["role_id"])
	roleName := params["role_name"].(string)
	zs := service.Role.Rename(roleId, roleName)
	ctx.JSON(zs)
}
func (c *RoleController) UpdateState(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	roleId := text.Conv.Int64(params["role_id"])
	roleState := text.Conv.Int(params["state"])
	zs := service.Role.UpdateRole(roleId, roleState)
	ctx.JSON(zs)
}
func (c *RoleController) Rename(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	roleId := text.Conv.Int64(params["role_id"])
	roleName := params["role_name"].(string)
	zs := service.Role.Rename(roleId, roleName)
	ctx.JSON(zs)
}

func (c *RoleController) SetPermission(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	roleId := text.Conv.Int64(params["role_id"])
	manageUserId := text.Conv.Int64(params["manage_user_id"])
	parentId := text.Conv.Int64(params["parent_role_id"])
	zs := service.Role.SetPermission(roleId, manageUserId, parentId)
	ctx.JSON(zs)
}

func (c *RoleController) SetRoleFunc(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	var reqData model.RoleFuncReq
	if err := mapstructure.Decode(params, &reqData); err != nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}

	zs := service.Role.SetRoleFunc(reqData)
	ctx.JSON(zs)
}

func (c *RoleController) GetRoleFunc(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	id := system.Text.Conv.Int64(params["role_id"])

	zs := service.Role.GetRoleFunc(id)
	ctx.JSON(zs)
}

func (c *RoleController) List(ctx iris.Context) {
	zs := service.Role.ListAll()
	ctx.JSON(zs)
}

func (c *RoleController) Find(ctx iris.Context) {
	data := ctx.Values().Get("data")
	condVal := ""
	if data != nil {
		params := data.(map[string]interface{})
		condVal = params["cond"].(string)
	}
	cond := "(a.role_name like ? or c.user_name like ?)"
	args := []any{"%" + condVal + "%", "%" + condVal + "%"}
	zs := service.Role.Find(cond, args)
	ctx.JSON(zs)
}

func (c *RoleController) ListFunc(ctx iris.Context) {
	data := ctx.Values().Get("data")
	cond := ""
	if data != nil {
		params := data.(map[string]interface{})
		if _, ok := params["cond"]; ok {
			cond = params["cond"].(string)
		}
	}
	zs := service.Role.ListFunc(cond)
	ctx.JSON(zs)
}
