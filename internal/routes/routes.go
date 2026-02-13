package routes

import (
	blogsHandler "github.com/Rajkumar-coderm/go-blog-backend/internal/handlers/blogs"
	commentsHandler "github.com/Rajkumar-coderm/go-blog-backend/internal/handlers/comments"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/handlers/profile"
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

	api.POST("/logout", middlewares.AuthMiddleware(), userhandler.LogoutUser)

	// Blog routes (authentication required)
	api.GET("/posts", middlewares.AuthMiddleware(), blogsHandler.GetAll)
	api.POST("/posts", middlewares.AuthMiddleware(), blogsHandler.CreatePost)
	api.PATCH("/posts/like", middlewares.AuthMiddleware(), blogsHandler.LikeDislikePost)
	api.DELETE("/posts", middlewares.AuthMiddleware(), blogsHandler.DeletePost)
	api.PATCH("/posts/bookmark", middlewares.AuthMiddleware(), blogsHandler.BookmarkPost)
	api.PATCH("/posts/save", middlewares.AuthMiddleware(), blogsHandler.SavedPost)

	// Comment routes (authentication required)
	api.POST("/posts/comment", middlewares.AuthMiddleware(), commentsHandler.CommentPost)
	api.GET("/posts/comment", middlewares.AuthMiddleware(), commentsHandler.GetAllPostComments)
	api.DELETE("/posts/comment", middlewares.AuthMiddleware(), commentsHandler.DeleteComment)

	// Profile routes (authentication required)
	api.GET("/profile", middlewares.AuthMiddleware(), profile.GetProfile)
	api.PATCH("/profile", middlewares.AuthMiddleware(), profile.UpdateProfile)
}
