package cache

import (
	"context"
	"fmt"

	"github.com/jhw66/myvideo_lab4/config"
	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis(cfg *config.AppConfig) {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}
	Rdb = redis.NewClient(options)

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		panic("Redis 连接失败: " + err.Error())
	}
}
