package profile

import (
	"net/http"

	"github.com/Rajkumar-coderm/go-blog-backend/internal/repositories/profile"
	"github.com/gin-gonic/gin"
)

func UpdateProfile(c *gin.Context) {

	err := profile.UpdateProfile(c)
	if err != nil {
		c.JSON(404, gin.H{
			"success": true,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
		"data":    nil,
	})
}

func GetProfile(c *gin.Context) {
	result, err := profile.GetProfile(c)
	if err != nil {
		c.JSON(401, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Success",
		"data":    result,
	})
}
