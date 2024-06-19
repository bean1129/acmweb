package server

import (
	"sync"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"

	"acmweb/system/config"
)

const (
	ModelHttp   = 1
	ModelTcp    = 2
	ModelWs     = 3
	ModelHandle = 4
	ModelDaemon = 5
)

var MyIris *iris.Application

func init() {
	if MyIris == nil {
		MyIris = iris.New()
		if config.CONFIG.Application.Debug {
			MyIris.Logger().SetLevel("debug")
		} else {
			MyIris.Logger().SetLevel("info")
		}
		MyIris.Use(logger.New())
		if config.CONFIG.Attachment.FilePath != "" {
			MyIris.HandleDir("/static/", config.CONFIG.Attachment.FilePath)
		}
	}
}

type Server struct {
	inst AppServer
}

func New(model int) *Server {
	wg := new(sync.WaitGroup)
	defer wg.Wait()

	c := new(Server)
	switch model {
	case ModelHttp:
		c.inst = NewHttpServ(wg)
	case ModelTcp:
		c.inst = NewTcpServ(wg)
	case ModelWs:
		c.inst = NewWsServ(wg)
	case ModelHandle:
		c.inst = NewHandleServ(wg)
	case ModelDaemon:
	default:
		c.inst = NewDaemonServ(wg)
	}

	return c
}

func (c *Server) Start() {
	c.inst.Start()
}

func (c *Server) Stop() {
	c.inst.Stop()
}

func (c *Server) Get() interface{} {
	return c.inst.Get()
}
