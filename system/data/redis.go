package data

import (
	"time"

	"github.com/kataras/iris/v12/sessions/sessiondb/redis"

	"acmweb/system/config"
)

type ZCache struct {
	Inst *redis.Database
}

func NewCache() *ZCache {
	redis := redis.New(initRedisConfig())
	if redis == nil {
		panic("Redis缓存数据库连接错误")
	}
	inst := new(ZCache)
	inst.Inst = redis
	return inst
}

// 初始化配置项
func initRedisConfig() redis.Config {

	// 默认配置
	cfg := redis.DefaultConfig()
	// 网络
	if config.CONFIG.Redis.Network != "" {
		cfg.Network = config.CONFIG.Redis.Network
	}
	// 地址
	if config.CONFIG.Redis.Addr != "" {
		cfg.Addr = config.CONFIG.Redis.Addr
	}
	// 超时时间
	if config.CONFIG.Redis.Timeout > 0 {
		cfg.Timeout = time.Duration(config.CONFIG.Redis.Timeout) * time.Second
	}
	// MaxActive
	if config.CONFIG.Redis.MaxActive > 0 {
		cfg.MaxActive = config.CONFIG.Redis.MaxActive
	}
	// 密码
	if config.CONFIG.Redis.Password != "" {
		cfg.Password = config.CONFIG.Redis.Password
	}
	// 数据库
	if config.CONFIG.Redis.Database != "" {
		cfg.Database = config.CONFIG.Redis.Database
	}
	// 前缀
	if config.CONFIG.Redis.Prefix != "" {
		cfg.Prefix = config.CONFIG.Redis.Prefix
	}
	return cfg
}

func (c *ZCache) Set() {
}
