package main

import (
	"backend-queue/config"
	"backend-queue/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.InitDB()
	config.InitRedis()

	r := gin.Default()
	routes.Route(r)
	r.Run(":8080")
}
