package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jhw66/myvideo_lab4/cache"
	"github.com/jhw66/myvideo_lab4/config"
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/router"
	"github.com/jhw66/myvideo_lab4/service"
	"github.com/jhw66/myvideo_lab4/utils"
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

	if err := utils.InitSnowflake(1); err != nil {
		log.Fatal(err)
	}

	if _, err := model.InitDB(cfg); err != nil {
		panic(err)
	}
	cache.InitRedis(cfg)

	service.WarmUpRankZSet()

	router.NewRouter(r)

	go service.StartVideoStatSync()

	r.Run(":" + cfg.RunPort)
}
