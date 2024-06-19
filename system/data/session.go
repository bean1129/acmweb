package data

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/sessions/sessiondb/redis"
)

type ZSession struct {
	Inst *sessions.Sessions
}

func NewSession(cookieName string, cache *redis.Database) *ZSession {
	iris.RegisterOnInterrupt(func() {
		cache.Close()
	})

	sess := sessions.New(sessions.Config{Cookie: cookieName})
	sess.UseDatabase(cache)
	inst := new(ZSession)
	inst.Inst = sess
	return inst
}
