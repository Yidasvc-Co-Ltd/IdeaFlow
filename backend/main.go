package main

import (
	"backend/app"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/api/backend", app.App_go)
	r.Run(":4504")
}
