package data

import (
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/sessions/sessiondb/redis"
)

var (
	DB      *ZDatabase
	Cache   *redis.Database
	Session *sessions.Sessions
)

func init() {
	DB = NewDB()
	Cache = NewCache().Inst
	Session = NewSession("acm.session", Cache).Inst
}
