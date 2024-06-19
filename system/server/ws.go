package server

import (
	"sync"
)

type WsServ struct {
	wg *sync.WaitGroup
}

func NewWsServ(wg *sync.WaitGroup) *WsServ {
	return &WsServ{wg}
}

func (c *WsServ) Start() {

}

func (c *WsServ) Stop() {

}

func (c *WsServ) Get() interface{} {
	return nil
}
