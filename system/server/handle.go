package server

import (
	"sync"
)

type HandleServ struct {
	wg *sync.WaitGroup
}

func NewHandleServ(wg *sync.WaitGroup) *DaemonServ {
	return &DaemonServ{wg}
}

func (c *HandleServ) Start() {

}

func (c *HandleServ) Stop() {

}

func (c *HandleServ) Get() interface{} {
	return nil
}
