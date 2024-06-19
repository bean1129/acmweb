package server

import (
	stdCtx "context"
	"sync"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/rs/cors"

	"acmweb/system/config"
	"acmweb/system/server/middleware"
)

type HttpServ struct {
	app *iris.Application
	wg  *sync.WaitGroup
}

func NewHttpServ(wg *sync.WaitGroup) *HttpServ {
	app := MyIris
	app.I18n.Load("./locales/*/*")
	app.I18n.SetDefault("zh-CN")

	// 跨域解决方案
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		Debug:            true,
	})
	app.WrapRouter(crs.ServeHTTP)
	app.Use(middleware.TranGen)
	//app.Use(middleware.TokenVerify)
	app.Use(middleware.ParseParams)
	app.Use(middleware.SecureVerify)
	return &HttpServ{app, wg}
}

func (c *HttpServ) Start() {
	iris.RegisterOnInterrupt(func() {
		c.wg.Add(1)
		defer c.wg.Done()

		ctx, cancel := stdCtx.WithTimeout(stdCtx.Background(), 20*time.Second)
		defer cancel()
		// 关闭所有主机
		c.app.Shutdown(ctx)
	})

	c.app.Run(iris.Addr(config.CONFIG.Application.Addr + ":" + config.CONFIG.Application.Port))
}

func (c *HttpServ) Stop() {
}

func (c *HttpServ) Get() interface{} {
	return c.app
}
