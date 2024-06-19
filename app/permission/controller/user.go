package controller

import (
	"acmweb/app/errors"
	"acmweb/app/permission/model"
	"acmweb/app/permission/service"
	"acmweb/app/permission/verification"
	"acmweb/system"
	"acmweb/system/common"
	"acmweb/system/text"
	"github.com/gookit/validate"
	"github.com/kataras/iris/v12"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

var User = new(UserController)

type UserController struct{}

func (c *UserController) Create(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	var u model.User
	if err := mapstructure.Decode(params, &u); err != nil {
		system.Log.Error("Invaild params,err:", err.Error())
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return

	}
	//校验参数
	uReq := &verification.UserReq{
		UserName: u.UserName,
		UserCode: u.UserCode,
		Password: u.Passwd,
		Phone:    u.Phone,
	}

	v := validate.Struct(uReq)
	if !v.Validate() {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), v.Errors.One()))
		return
	}
	//密码加密
	u.Passwd, _ = system.Secret.MD5.Password(u.Passwd + system.Config.Application.SecretKey + u.UserCode)
	//创建
	zs, err := service.User.CreateUser(u)
	if err != nil {
		zs.Code = errors.ErrUserExists.ErrorCode()
		zs.Msg = ctx.Tr(err.Error())
		system.Log.Error("Create user failed,errMsg:", zs.Msg)
	}
	ctx.JSON(zs)
}

func (c *UserController) Delete(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	ids := make([]int64, 1)

	//admin不允许操作

	//是否批量
	v := reflect.ValueOf(params["user_id"])
	if v.Kind() == reflect.Slice {
		idsparams := params["user_id"].([]interface{})
		for _, v := range idsparams {
			ids = append(ids, int64(v.(float64)))
		}
	} else {
		ids = append(ids, system.Text.Conv.Int64(params["user_id"]))
	}

	zs, err := service.User.DeleteUser(ids)
	if err != nil {
		system.Log.Error("Delete user failed,errMsg:", zs.Msg)
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}

func (c *UserController) ModifyUser(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}

	params := data.(map[string]interface{})
	var u model.User
	if err := mapstructure.Decode(params, &u); err != nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
	}
	//校验参数
	phoneReq := &verification.PhoneReq{
		Phone: u.Phone,
	}

	v := validate.Struct(phoneReq)
	if !v.Validate() {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), v.Errors.One()))
		return
	}

	//修改信息
	zs, err := service.User.ModifyUser(u)
	if err != nil {
		system.Log.Error("Update user failed,errMsg:", zs.Msg)
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}

func (c *UserController) RenameUser(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	m := model.User{UserId: text.Conv.Int64(params["user_id"]), UserName: params["user_name"].(string)}
	zs, err := service.User.Rename(m)
	if err != nil {
		system.Log.Error("Rename user failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
func (c *UserController) SetUserRole(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	m := model.User{UserId: text.Conv.Int64(params["user_id"]), RuleId: text.Conv.Int64(params["role_id"])}
	zs, err := service.User.UpdateRole(m)
	if err != nil {
		system.Log.Error("Rename user failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
func (c *UserController) ChangeRole(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	m := model.User{UserId: text.Conv.Int64(params["user_id"]), RuleId: text.Conv.Int64(params["role_id"])}
	zs, err := service.User.ChangeRole(m)
	if err != nil {
		system.Log.Error("ChangeRole user failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}

func (c *UserController) SetPasswd(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	m := model.User{UserId: text.Conv.Int64(params["user_id"]), Passwd: params["password"].(string)}
	zs, err := service.User.SetPassword(m)
	if err != nil {
		system.Log.Error("SetPasswd user failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}

func (c *UserController) UpdateState(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	m := model.User{UserId: text.Conv.Int64(params["user_id"]), State: text.Conv.Int(params["state"])}
	zs, err := service.User.ChangeState(m)
	if err != nil {
		system.Log.Error("Dis user failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}

func (c *UserController) Disable(ctx iris.Context) {
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	m := model.User{UserId: text.Conv.Int64(params["user_id"]), State: text.Conv.Int(params["State"])}
	zs, err := service.User.ChangeState(m)
	if err != nil {
		system.Log.Error("Disable user failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}

func (c *UserController) GetUsersInfo(ctx iris.Context) {
	zs := common.NewResult()
	data := ctx.Values().Get("data")
	if data == nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	params := data.(map[string]interface{})
	u := &model.UserSearch{State: -1, Vague: true}
	if err := mapstructure.Decode(params, u); err != nil {
		ctx.JSON(common.NewErrResult(errors.ErrParamMiss.ErrorCode(), errors.ErrParamMiss.Error()))
		return
	}
	zs, err := service.User.List(u)
	if err != nil {
		system.Log.Error("getPageUsers  failed,errMsg:", err.Error())
		zs.Msg = ctx.Tr(err.Error())
	}
	ctx.JSON(zs)
}
