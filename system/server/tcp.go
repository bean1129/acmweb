package server

import (
	"sync"
)

type TcpServ struct {
	wg *sync.WaitGroup
}

func NewTcpServ(wg *sync.WaitGroup) *TcpServ {
	return &TcpServ{wg}
}

func (c *TcpServ) Start() {

}

func (c *TcpServ) Stop() {

}

func (c *TcpServ) Get() interface{} {
	return nil
}
