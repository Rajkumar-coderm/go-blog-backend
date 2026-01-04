package blogsHandler

import (
	"net/http"

	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/repositories/blogs"
	"github.com/gin-gonic/gin"
)

// Create a new post
func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDStr, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	result, err := blogs.CreatePost(&post, userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post created", "id": result.InsertedID})
}

// GetAll - Fetch all posts with author details
func GetAll(c *gin.Context) {
	var request models.GetRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters", "message": err.Error()})
		return
	}
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDStr, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}
	if request.Skip < 0 {
		request.Skip = 0
	}
	if request.Limit < 1 {
		request.Limit = 40
	}
	posts, totalCount, err := blogs.GetAll(c, &request, userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch posts", "message": err.Error()})
		return
	}

	APIResponse := models.APIResponse{
		Status:     "success",
		Data:       posts,
		TotalCount: totalCount,
	}

	c.JSON(http.StatusOK, APIResponse)
}

func LikeDislikePost(c *gin.Context) {
	finalResponse := models.CommonGetResponse{}
	err := blogs.LikeDislikePost(c)
	if err != nil {
		finalResponse.Message = err.Error()
		finalResponse.Success = false
		finalResponse.Data = nil
		c.JSON(http.StatusBadRequest, finalResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Request successfully completed",
		"message": "Success",
		"data":    nil})
}

func BookmarkPost(c *gin.Context) {
	finalResponse := models.CommonGetResponse{}
	err := blogs.BookmarkPost(c)
	if err != nil {
		finalResponse.Message = err.Error()
		finalResponse.Success = false
		finalResponse.Data = nil
		c.JSON(http.StatusBadRequest, finalResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Request successfully completed",
		"message": "Success",
		"data":    nil})
}

func DeletePost(c *gin.Context) {
	finalResponse := models.CommonGetResponse{}
	err := blogs.DeletePost(c)
	if err != nil {
		finalResponse.Message = err.Error()
		finalResponse.Success = false
		finalResponse.Data = nil
		c.JSON(http.StatusBadRequest, finalResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Request successfully completed",
		"message": "Success",
		"data":    nil})
}

func SavedPost(c *gin.Context) {
	blogs.SavePost(c)
}
