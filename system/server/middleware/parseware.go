package middleware

import (
	"encoding/json"
	"github.com/kataras/iris/v12"

	"acmweb/system/web"
)

func _parseJson(ctx iris.Context) bool {
	var params map[string]interface{}
	if err := ctx.ReadJSON(&params); err != nil {
		ctx.JSON(web.JsonResult{
			Code: 400,
			Msg:  "BODY格式不正确",
		})
		return false
	}
	if params["time"] == nil || params["sign"] == nil {
		ctx.JSON(web.JsonResult{
			Code: 400,
			Msg:  "缺少认证信息",
		})
		return false
	}
	if params["data"] == nil {
		ctx.JSON(web.JsonResult{
			Code: 400,
			Msg:  "缺少数据对象",
		})
		return false
	}
	ctx.Values().Set("time", params["time"])
	ctx.Values().Set("sign", params["sign"])
	ctx.Values().Set("data", params["data"])
	return true
}

func _parseForm(ctx iris.Context) bool {
	var params map[string]interface{}
	if err := ctx.ReadForm(&params); err != nil {
		ctx.JSON(web.JsonResult{
			Code: 400,
			Msg:  "参数有误解析失败",
		})
		return false
	}
	ctx.Values().Set("data", params)
	return true
}

func _parseMultiForm(ctx iris.Context) bool {
	var params map[string]interface{}
	paramstr := ctx.FormValue("app-data")
	if len(paramstr) > 0 {
		json.Unmarshal([]byte(paramstr), &params)
		ctx.Values().Set("data", params["data"].(map[string]interface{}))
		ctx.Values().Set("time", params["time"])
		ctx.Values().Set("sign", params["sign"])
	}
	return true
}

func ParseParams(ctx iris.Context) {
	tranId := ctx.Values().Get("tranId")
	ctx.Application().Logger().Infof("[%v] Parse request body and convert map struct", tranId)
	ct := ctx.GetHeader("Content-Type")
	if ct == "application/json" {
		if _parseJson(ctx) == false {
			return
		}
	} else if ct == "application/x-www-form-urlencoded" {
		if _parseForm(ctx) == false {
			return
		}
	} else if len(ct) >= 19 && ct[0:19] == "multipart/form-data" {
		if _parseMultiForm(ctx) == false {
			return
		}
	}
	ctx.Next()
}
