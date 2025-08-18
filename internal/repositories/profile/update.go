package profile

import (
	"context"
	"errors"
	"time"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateProfile(c *gin.Context) error {
	col := config.DB.Collection("users")

	// Extract userID from context
	userID, ok := c.Get("userID")
	if !ok {
		return errors.New("user ID not found")
	}

	_id, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		return errors.New("invalid user ID format")
	}

	// Fetch existing user data from DB
	var existingUser models.User
	err = col.FindOne(context.TODO(), bson.M{"_id": _id}).Decode(&existingUser)
	if err != nil {
		return err
	}

	// Bind request body to updateUser
	var updateUser map[string]interface{}
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		return err
	}

	// Ensure sensitive fields are not updated
	delete(updateUser, "password")
	delete(updateUser, "role")
	delete(updateUser, "emailVerified")
	delete(updateUser, "active")
	delete(updateUser, "createdAt")
	delete(updateUser, "updatedAt")
	delete(updateUser, "lastLogin")

	// Remove fields with null or empty values
	for key, value := range updateUser {
		if value == nil || value == "" {
			delete(updateUser, key)
		}
	}

	updateUser["updatedAt"] = time.Now()
	update := bson.M{"$set": updateUser}
	result := col.FindOneAndUpdate(context.TODO(), bson.M{"_id": _id}, update)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}
