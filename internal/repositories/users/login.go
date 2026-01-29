package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/auth"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/repositories/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

func LoginUser(c *gin.Context, request models.LoginRequest) (*models.TokenModel, error) {
	col := config.DB.Collection("users")
	var user models.User
	var filter bson.M

	// Set filter based on login type
	switch request.LoginType {
	case "email":
		filter = bson.M{"email": request.Email, "registration_type": "email"}
	case "phone":
		filter = bson.M{"phone": request.Phone, "country_code": request.CountryCode, "phone_iso_code": request.PhoneIsoCode, "registration_type": "phone"}
	case "google":
		filter = bson.M{"google_id": request.GoogleID, "registration_type": "google"}
	default:
		return nil, errors.New("invalid login type")
	}

	// Log request for debugging
	fmt.Println("Login filter:", filter)

	// Try to find the user
	err := col.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Optional fallback: find by email alone (in case registration_type is missing)
			if request.LoginType == "email" {
				fmt.Println("Fallback: trying without registration_type")
				err = col.FindOne(context.TODO(), bson.M{"email": request.Email}).Decode(&user)
				if err == mongo.ErrNoDocuments {
					return nil, ErrUserNotFound
				} else if err != nil {
					return nil, err
				}
			} else {
				return nil, ErrUserNotFound
			}
		} else {
			return nil, err
		}
	}

	// Verify credentials
	switch request.LoginType {
	case "email", "phone":
		if user.Password == "" {
			return nil, ErrInvalidCredentials
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
		if err != nil {
			return nil, ErrInvalidCredentials
		}
	case "google":
		if user.GoogleID != request.GoogleID {
			return nil, ErrInvalidCredentials
		}
	}

	// Generate tokens
	token, err := auth.GenerateJWT(user.ID.Hex())
	if err != nil {
		return nil, err
	}
	refreshToken, err := auth.GenerateRefreshToken(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	// Create session
	deviceInfo := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()
	sessionExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	_, err = sessions.CreateSession(user.ID, refreshToken, deviceInfo, ipAddress, sessionExpiresAt)
	if err != nil {
		return nil, err
	}

	// Update user's login timestamp
	_, err = col.UpdateOne(
		context.TODO(),
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"active": true, "last_login": time.Now()}},
	)
	if err != nil {
		return nil, err
	}

	// Prepare response
	var userToken models.TokenModel
	userToken.ID = user.ID.Hex()
	userToken.Name = user.FirstName + " " + user.LastName
	userToken.Username = user.Username
	userToken.Email = user.Email
	userToken.Role = user.Role
	userToken.EmailVerified = user.EmailVerified
	userToken.CreatedAt = user.CreatedAt
	userToken.UpdatedAt = user.UpdatedAt
	userToken.Active = user.Active
	userToken.Token = map[string]interface{}{
		"token":                 token,
		"type":                  "Bearer",
		"expiresIn":             15 * time.Minute,
		"expiresAt":             time.Now().Add(15 * time.Minute),
		"refreshToken":          refreshToken,
		"refreshTokenExpiresIn": 7 * 24 * time.Hour,
		"refreshTokenExpiresAt": time.Now().Add(7 * 24 * time.Hour),
	}

	return &userToken, nil
}
