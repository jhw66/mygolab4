package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jhw66/myvideo_lab4/cache"
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/router"
	"github.com/jhw66/myvideo_lab4/service"
)

func main() {
	r := gin.Default()
	r.Static("/static", "./static")

	if _, err := model.InitDB(); err != nil {
		panic(err)
	}
	cache.InitRedis()

	router.NewRouter(r)
	go service.SyncFavoirte()
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		for range ticker.C {
			service.SyncFavoriteCount()
			service.SyncCommentCount()
		}
	}()

	r.Run(":80")
}
