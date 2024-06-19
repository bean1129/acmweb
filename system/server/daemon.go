package server

import (
	"sync"
)

type DaemonServ struct {
	wg *sync.WaitGroup
}

func NewDaemonServ(wg *sync.WaitGroup) *DaemonServ {
	return &DaemonServ{wg}
}

func (c *DaemonServ) Start() {

}

func (c *DaemonServ) Stop() {

}

func (c *DaemonServ) Get() interface{} {
	return nil
}
