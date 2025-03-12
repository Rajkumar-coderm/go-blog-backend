package routes

import (
	blogsHandler "github.com/Rajkumar-coderm/go-blog-backend/internal/handlers/blogs"
	userhandler "github.com/Rajkumar-coderm/go-blog-backend/internal/handlers/users"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// Add CORS middleware
	r.Use(middlewares.CORSMiddleware())

	// Global OPTIONS handler
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204)
	})

	api := r.Group("/v1")

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Page not found"})
	})

	// User routes (no authentication required)
	api.POST("/register", userhandler.RegisterUser)
	api.POST("/login", userhandler.LoginUser)
	api.GET("/validate-username", userhandler.ValidateUserName)

	// Blog routes (authentication required)
	api.Use(middlewares.AuthMiddleware())
	api.GET("/posts", blogsHandler.GetAll)
	api.POST("/posts", blogsHandler.CreatePost)
	api.PATCH("/posts/like", blogsHandler.LikeDislikePost)
	api.PATCH("/posts/bookmark", blogsHandler.BookmarkPost)
	api.POST("/posts/comment", blogsHandler.CommentPost)
	api.GET("/posts/comment", blogsHandler.GetAllPostComments)
}
