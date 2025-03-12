package main

import (
	"os"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()
	r := gin.Default()
	routes.RegisterRoutes(r)
	port := os.Getenv("PORT")
	if port != "" {
		r.Run(":" + port)
	} else {
		r.Run(":8080")
	}
}
