package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	options := &redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       10,
	}
	Rdb = redis.NewClient(options)

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		panic("Redis 连接失败: " + err.Error())
	}
}
