package middleware

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/kataras/iris/v12"

	"acmweb/constants"
	"acmweb/system/config"
	"acmweb/system/secret"
	"acmweb/system/text"
	"acmweb/system/web"
)

func SecureVerify(ctx iris.Context) {
	tranId := ctx.Values().Get("tranId")
	ctx.Application().Logger().Infof("[%v] Verify request body is valid", tranId)
	if config.CONFIG.Application.Debug {
		ctx.Next()
		return
	}

	var (
		params map[string]interface{}
		_time  string
		sign   string
	)

	params = ctx.Values().Get("data").(map[string]interface{})
	_time = ctx.Values().GetString("time")
	sign = ctx.Values().GetString("sign")

	t1 := time.Now().Unix()
	t2, _ := strconv.ParseInt(_time, 10, 64)
	// 超过5分钟认为接口请求过期
	if t1-t2 > 300 {
		ctx.Application().Logger().Infof("[%v] Time expire t1:%v t2:%v", tranId, t1, t2)
		ctx.JSON(web.JsonResult{
			Code: constants.ErrProtoTimeExpire,
			Msg:  "请求时间戳过期",
		})
		return
	}

	marshal, err := json.Marshal(params)
	if err != nil {
		ctx.JSON(web.JsonResult{
			Code: constants.ErrProtoSignInvalid,
			Msg:  "不正确的JSON数据",
		})
	}
	bodyData := string(marshal)
	originStr := text.String.SubStr(config.CONFIG.Application.SecretKey, 0, 4)
	originStr += text.String.SubStr(bodyData, 0, 10)
	originStr += _time
	originStr += text.String.SubStr(config.CONFIG.Application.SecretKey, 4)
	originStr += text.String.SubStr(bodyData, 10)

	s1, err := secret.MD5.MD5(originStr)
	if err != nil || s1 != sign {
		ctx.Application().Logger().Infof("[%v] Sign invalid, the correct value: %v", tranId, s1)
		ctx.JSON(web.JsonResult{
			Code: constants.ErrProtoSignInvalid,
			Msg:  "请求签名不正确",
		})
		return
	}
	ctx.Next()
}
