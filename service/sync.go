package service

import (
	"time"
)

func StartVideoStatSync() {
	ticker := time.NewTicker(SyncInterval)
	defer ticker.Stop()
	for range ticker.C {
		SyncDirtyToMySQL()
	}
}
