package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/router"
)

func main() {
	r := gin.Default()
	r.Static("/static", "./static")

	if _, err := model.InitDB(); err != nil {
		panic(err)
	}

	router.NewRouter(r)

	r.Run(":80")
}
