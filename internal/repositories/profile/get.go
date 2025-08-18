package profile

import (
	"context"
	"errors"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetProfile(c *gin.Context) (interface{}, error) {
	col := config.DB.Collection("users")

	// Get user ID from context or query parameter
	userID, ok := c.Get("userID")
	if queryID := c.Query("id"); queryID != "" {
		userID = queryID
	}

	if !ok {
		return nil, errors.New("user ID not found")
	}

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Fetch user from database
	var user models.User
	if err := col.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user); err != nil {
		return nil, err
	}

	// Return user data as map
	return map[string]interface{}{
		"id":            user.ID,
		"firstName":     user.FirstName,
		"lastName":      user.LastName,
		"email":         user.Email,
		"username":      user.Username,
		"role":          user.Role,
		"emailVerified": user.EmailVerified,
		"lastLogin":     user.LastLogin,
		"contact":       user.Contact,
		"blog":          user.Blog,
		"activity":      user.Activity,
		"birthday":      user.Birthday,
		"avatar":        user.Avatar,
		"location":      user.Location,
		"notifications": user.Notifications,
		"bio":           user.Bio,
		"active":        user.Active,
		"createdAt":     user.CreatedAt,
		"updatedAt":     user.UpdatedAt,
		"phone":         user.Phone,
		"countryCode": user.CoutryCode,
		"phoneIsoCode": user.PhoneIsoCode,
	}, nil
}
