package common

import (
	"acmweb/system/config"
)

var (
	Utils *ZUtils
	UUID  *ZUUID
)

func init() {
	Utils = NewUtils()
	UUID = NewUUID(int64(config.CONFIG.Application.ServerId))
}
