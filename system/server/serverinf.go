package server

import (
	"github.com/kataras/iris/v12/context"
)

type RouterAdapter struct {
	Party  string
	Router map[string]context.Handler
}

type AppServer interface {
	Start()
	Stop()
	Get() interface{}
}
