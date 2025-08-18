package main

import (
	"os"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/routes"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()
	r := gin.Default()
	routes.RegisterRoutes(r)

	// Start socket server in a goroutine so it doesn't block the main HTTP server
	go services.InitSocketServer()

	port := os.Getenv("PORT")
	if port != "" {
		r.Run(":" + port)
	} else {
		r.Run(":8080")
	}
}
