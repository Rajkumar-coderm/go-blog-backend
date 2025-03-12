package users

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/auth"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context, user *models.User) (*models.TokenModel, error) {
	col := config.DB.Collection("users")
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	fmt.Println(user)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if user.Email == "" {
		return nil, errors.New("email cannot be empty")
	}

	var existingUser models.User
	err = col.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		return nil, errors.New("email already exists")
	} else if err != mongo.ErrNoDocuments {
		return nil, err
	}

	if user.Username == "" {
		user.Username = strings.Split(user.Email, "@")[0]
	}

	err = col.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		return nil, errors.New("username already exists")
	} else if err != mongo.ErrNoDocuments {
		return nil, err
	}
	user.Password = string(hashedPassword)
	_, error := col.InsertOne(context.TODO(), user)
	if error != nil {
		return nil, error
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
	var userToken models.TokenModel
	userToken.ID = user.ID.Hex()
	userToken.Name = user.FirstName + " " + user.LastName
	userToken.Username = user.Username
	userToken.Email = user.Email
	userToken.Role = user.Role
	userToken.EmailVerified = user.EmailVerified
	userToken.Contact = user.Contact.Phone
	userToken.CreatedAt = user.CreatedAt
	userToken.UpdatedAt = user.UpdatedAt
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
