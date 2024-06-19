package middleware

import (
	"time"

	"github.com/kataras/iris/v12"

	"acmweb/system/common"
	"acmweb/system/config"
	"acmweb/system/data"
	"acmweb/system/text"
	"acmweb/system/web"
)

// TokenVerify 登录验证中间件
func TokenVerify(ctx iris.Context) {
	// 生成交易流水号
	tranId := ctx.Values().Get("tranId")
	ctx.Application().Logger().Infof("[%v] Verify token is valid", tranId)

	// 放行设置
	urlItem := []string{"/captcha", "/login", "/connect"}
	if !common.Utils.InStringArray(ctx.Path(), urlItem) {
		// 检查是否用户已经认证过
		if auth, _ := data.Session.Start(ctx).GetBoolean("authenticated"); !auth {
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(web.JsonResult{
				Code: 400,
				Msg:  "未找到登录信息",
			})
			return
		}

		// 从请求头中获取Token
		token := ctx.GetHeader("Authorization")
		// 字符串替换
		token = text.String.Replace(token, "Bearer ", "")
		claim, err := web.Token.Parse(token)
		if err != nil {
			ctx.Application().Logger().Warnf("[%v] Token parse error, %v", tranId, err)
			ctx.JSON(web.JsonResult{
				Code: 401,
				Msg:  "Token已过期",
			})
			return
		} else if time.Now().Unix() > claim.ExpiresAt {
			ctx.Application().Logger().Warnf("[%v] Token timeout, please login in again", tranId)
			ctx.JSON(web.JsonResult{
				Code: 401,
				Msg:  "时间超时",
			})
			return
		}
		session := data.Session.Start(ctx)
		username := session.Get("userid")
		if username != claim.Username {
			ctx.Application().Logger().Infof("[%v] Token invalid, userid: %s", tranId, username)
			ctx.JSON(web.JsonResult{
				Code: 401,
				Msg:  "Token无效",
			})
			return
		}
		// 校验通过，更新时间戳
		nowTime := time.Now()
		expireTime := nowTime.Add(time.Duration(config.CONFIG.Application.ExpireTime) * time.Hour)
		claim.StandardClaims.ExpiresAt = expireTime.Unix()
	}
	// 前置中间件
	// ctx.Application().Logger().Infof("Runs before %s", ctx.Path())
	ctx.Next()
}
