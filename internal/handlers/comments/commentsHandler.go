package comments

import (
	"net/http"

	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/repositories/comments"
	"github.com/gin-gonic/gin"
)

func CommentPost(c *gin.Context) {
	finalResponse := models.CommonGetResponse{}
	err := comments.CommentPost(c)
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

func GetAllPostComments(c *gin.Context) {
	var request models.GetRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters", "message": err.Error()})
		return
	}
	comments, totalCount, err := comments.GetAllPostComments(c, &request)
	if err != nil {
		APIResponse := models.APIResponse{
			Status:     "error",
			Data:       nil,
			TotalCount: 0,
			Message:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, APIResponse)
		return
	}
	APIResponse := models.APIResponse{
		Status:     "success",
		Data:       comments,
		TotalCount: totalCount,
		Message:    "Success",
	}

	c.JSON(http.StatusOK, APIResponse)
}

func DeleteComment(c *gin.Context) {
	err := comments.DeleteComment(c)
	if err != nil {
		APIResponse := models.APIResponse{
			Status:     "error",
			Data:       nil,
			TotalCount: 0,
			Message:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, APIResponse)
		return
	}
	APIResponse := models.APIResponse{
		Status:  "success",
		Data:    nil,
		Message: "Success",
	}
	c.JSON(http.StatusOK, APIResponse)
}
