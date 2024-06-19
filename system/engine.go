package system

import (
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/sessions/sessiondb/redis"

	"acmweb/system/common"
	"acmweb/system/config"
	"acmweb/system/data"
	"acmweb/system/secret"
	"acmweb/system/server"
	"acmweb/system/text"
	"acmweb/system/web"
)

type _data struct {
	DB      *data.ZDatabase
	Cache   *redis.Database
	Session *sessions.Sessions
}

type _secret struct {
	MD5 *secret.ZMD5
	DES *secret.ZDES
}

type _text struct {
	Conv   *text.ZConv
	Regex  *text.ZRegex
	File   *text.ZFile
	String *text.ZString
}

type _web struct {
	Token *web.ZToken
}

type _com struct {
	Utils *common.ZUtils
	UUID  *common.ZUUID
}

var (
	Data   = new(_data)
	Secret = new(_secret)
	Text   = new(_text)
	Web    = new(_web)
	Common = new(_com)
	Config *config.ResourceConfig
	Log    *golog.Logger
)

func init() {
	Data.DB = data.DB
	Data.Cache = data.Cache
	Data.Session = data.Session

	Secret.MD5 = secret.MD5
	Secret.DES = secret.DES

	Text.Conv = text.Conv
	Text.Regex = text.Regex
	Text.File = text.File
	Text.String = text.String

	Web.Token = web.Token

	Common.Utils = common.Utils
	Common.UUID = common.UUID

	Config = &config.CONFIG

	Log = server.MyIris.Logger()
	Log.SetLevel(Config.Application.LogLevel)
}
