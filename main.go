package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"wishwall/app/midwares"
	"wishwall/config/config"
	"wishwall/config/database"
	"wishwall/config/router"
	"wishwall/config/session"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	database.Init()

	r := gin.Default()
	r.Use(midwares.Cors())
	r.Use(midwares.ErrHandler())
	r.NoMethod(midwares.HandleNotFound)
	r.NoRoute(midwares.HandleNotFound)

	session.Init(r)
	router.Init(r)

	err := r.Run(":" + config.Config.GetString("router.port"))
	if err != nil {
		log.Fatal("ServerStartFailed", err)
	}
}
