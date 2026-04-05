package cache

import (
	"time"
)

func TryWarmupLock(key string, ttl time.Duration) bool {
	ok, err := Rdb.SetNX(Ctx, key, 1, ttl).Result()
	if err != nil {
		return false
	}
	return ok
}
