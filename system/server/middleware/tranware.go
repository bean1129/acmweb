package middleware

import (
	"github.com/kataras/iris/v12"

	"acmweb/system/common"
)

func TranGen(ctx iris.Context) {
	// 生成交易流水号
	tranId := common.UUID.NextVal()
	ctx.Values().Set("tranId", tranId)
	ctx.Next()
}
