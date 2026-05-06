package cachemodel

import (
	"context"
	"fmt"

	"github.com/jhw66/myvideo_lab4/internal/config"
	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis(cfg *config.Config) (*redis.Client, error) {
	options := &redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	}
	rdb := redis.NewClient(options)
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("ping redis failed:%w", err)
	}

	Rdb = rdb
	return rdb, nil
}
