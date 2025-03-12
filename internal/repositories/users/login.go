package users

import (
	"context"
	"errors"
	"time"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/auth"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
)

func LoginUser(email, password string) (*models.TokenModel, error) {
	col := config.DB.Collection("users")
	// Find the user by email
	var user models.User
	var userToken models.TokenModel
	err := col.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT auth token
	token, err := auth.GenerateJWT(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	/// Generate refresh token
	refreshToken, referr := auth.GenerateRefreshToken(user.ID.Hex())
	if referr != nil {
		return nil, referr
	}
	col.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"active": true, "lastLogin": time.Now()}})
	user.Token = token
	user.RefreshToken = refreshToken
	userToken.ID = user.ID.Hex()
	userToken.Email = user.Email
	userToken.Role = user.Role
	userToken.EmailVerified = user.EmailVerified
	userToken.Contact = user.Contact.Phone
	userToken.CreatedAt = user.CreatedAt
	userToken.UpdatedAt = user.UpdatedAt
	userToken.Active = user.Active
	userToken.Name = user.FirstName + " " + user.LastName
	userToken.Username = user.Username
	userToken.Token = map[string]interface{}{
		"token":                 token,
		"type":                  "Bearer",
		"expiresIn":             24 * time.Hour,
		"expiresAt":             time.Now().Add(24 * time.Hour),
		"refreshToken":          refreshToken,
		"refreshTokenExpiresIn": 7 * 24 * time.Hour,
		"refreshTokenExpiresAt": time.Now().Add(7 * 24 * time.Hour),
	}
	return &userToken, nil
}
