package main

import (
	"backend-queue/config"
	"backend-queue/routes"
	"backend-queue/services"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.InitDB()
	config.InitRedis()

	go services.Hub.Run()

	r := gin.Default()
	routes.Route(r)
	r.Run(":8080")
}
