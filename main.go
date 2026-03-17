package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jhw66/myvideo_lab4/cache"
	"github.com/jhw66/myvideo_lab4/config"
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/router"
	"github.com/jhw66/myvideo_lab4/service"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()
	r.Static("/static", "./static")

	err := godotenv.Load()
	if err != nil {
		log.Panicln(".env file not found")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := model.InitDB(cfg); err != nil {
		panic(err)
	}
	cache.InitRedis(cfg)

	router.NewRouter(r)
	go service.SyncFavoirte()
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		for range ticker.C {
			service.SyncFavoriteCount()
			service.SyncCommentCount()
		}
	}()

	r.Run(":" + cfg.RunPort)
}
