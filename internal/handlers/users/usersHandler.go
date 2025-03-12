package usershandeler

import (
	"fmt"
	"net/http"

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
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	finalResponse := models.CommonGetResponse{}

	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println(err)
		finalResponse.Message = "Something went wrong please try again " + err.Error()
		finalResponse.Success = false
		finalResponse.Data = nil
		c.JSON(http.StatusBadRequest, finalResponse)
		return
	}
	res, err := users.LoginUser(request.Email, request.Password)
	if err != nil {
		if err == users.ErrUserNotFound || err == users.ErrInvalidCredentials {
			finalResponse.Message = err.Error()
			finalResponse.Success = false
			finalResponse.Data = nil
			c.JSON(http.StatusUnauthorized, finalResponse)
			return
		} else {
			finalResponse.Message = "Something went wrong"
			finalResponse.Success = false
			finalResponse.Data = nil
			c.JSON(http.StatusInternalServerError, finalResponse)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data":    res,
	})
}
