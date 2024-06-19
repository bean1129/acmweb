package main

import (
	"github.com/kataras/iris/v12"

	"acmweb/router"
	"acmweb/system"
	"acmweb/system/server"
)

func main() {
	system.Log.Info(`----------------------------------------------------------------`)
	system.Log.Info(`      =       =====   |\      /|     ===== `)
	system.Log.Info(`     = =    |         | \    / |   |       `)
	system.Log.Info(`    =   =   |         |  \  /  |     ====  `)
	system.Log.Info(`   = === =  |         |   \/   |         | `)
	system.Log.Info(`  =       =   =====   |        |    =====  `)
	system.Log.Info(`----------------------------------------------------------------`)
	system.Log.Info("Welcome to use ppplication authorization management system")
	system.Log.Info("Prepare server capacity")
	system.Log.Info("Start data service server")
	serv := server.New(1)
	app := serv.Get().(*iris.Application)
	system.Log.Info("Register application router")
	router.RegisterRouter(app)
	serv.Start()
}
