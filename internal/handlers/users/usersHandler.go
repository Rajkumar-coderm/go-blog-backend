package usershandeler

import (
	"log"
	"net/http"

	"github.com/Rajkumar-coderm/go-blog-backend/internal/auth"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/repositories/users"
	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := users.RegisterUser(c, &user)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User registered successfully",
		"data":    result,
	})
}

func ValidateUserName(c *gin.Context) {
	_, err := users.VerifyUserName(c)
	if err != nil {
		c.JSON(404, gin.H{"status": "Reqest Failed",
			"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Request successfully completed",
		"message": "Success"})
}

func LoginUser(c *gin.Context) {
	var request models.LoginRequest
	var response models.CommonGetResponse

	// Bind and validate request
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Login: Invalid request body: %v\n", err)
		response = models.CommonGetResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Attempt login
	res, err := users.LoginUser(c, request)
	if err != nil {
		log.Printf("Login error: %v\n", err)
		status := http.StatusInternalServerError
		message := "Something went wrong. Please try again."

		// Handle known user errors
		if err == users.ErrUserNotFound || err == users.ErrInvalidCredentials {
			status = http.StatusUnauthorized
			message = err.Error()
		}

		response = models.CommonGetResponse{
			Success: false,
			Message: message,
		}
		c.JSON(status, response)
		return
	}

	// Success response
	c.JSON(http.StatusOK, models.CommonGetResponse{
		Success: true,
		Message: "Login successful",
		Data:    res,
	})
}

func LogoutUser(c *gin.Context) {
	var request models.LogoutRequest
	var response models.CommonGetResponse

	// Bind and validate request
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Logout: Invalid request body: %v\n", err)
		response = models.CommonGetResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get user ID from JWT token
	userID, err := auth.ExtractUserIDFromToken(c)
	if err != nil {
		log.Printf("Logout: Unable to extract user ID: %v\n", err)
		c.JSON(http.StatusUnauthorized, models.CommonGetResponse{
			Success: false,
			Message: "Unauthorized: " + err.Error(),
		})
		return
	}

	// Perform logout
	err = users.LogoutUser(c, userID, request.RefreshToken, request.LogoutAll)
	if err != nil {
		log.Printf("Logout error: %v\n", err)
		response = models.CommonGetResponse{
			Success: false,
			Message: "Error during logout: " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Success response
	logoutType := "current device"
	if request.LogoutAll {
		logoutType = "all devices"
	}

	c.JSON(http.StatusOK, models.CommonGetResponse{
		Success: true,
		Message: "Logged out from " + logoutType + " successfully",
	})
}
