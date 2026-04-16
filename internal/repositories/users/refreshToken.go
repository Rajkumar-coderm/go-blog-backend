package users

import (
	"context"
	"errors"
	"time"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/auth"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/repositories/sessions"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrRefreshTokenExpired = errors.New("refresh token expired or invalid")
)

func RefreshToken(c *gin.Context, request models.RefreshTokenRequest) (*models.TokenModel, error) {
	claims, err := auth.ValidateJWT(request.RefreshToken)
	if err != nil {
		return nil, ErrRefreshTokenExpired
	}

	if claims.Type != "refresh" {
		return nil, ErrInvalidRefreshToken
	}

	session, err := sessions.GetSessionByToken(request.RefreshToken)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrRefreshTokenExpired
	}

	userID, err := primitive.ObjectIDFromHex(claims.ID)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	if session.UserID != userID {
		return nil, ErrInvalidRefreshToken
	}

	col := config.DB.Collection("users")
	var user models.User
	if err := col.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	newAccessToken, err := auth.GenerateJWT(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := auth.GenerateRefreshToken(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	deviceInfo := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()
	if ipAddress == "::1" {
		ipAddress = "127.0.0.1"
	}
	newSessionExpiresAt := time.Now().Add(utils.RefreshTokenExpiryDuration)

	if err := sessions.DeleteSessionByToken(request.RefreshToken); err != nil {
		return nil, err
	}

	_, err = sessions.CreateSession(user.ID, newRefreshToken, deviceInfo, ipAddress, newSessionExpiresAt)
	if err != nil {
		return nil, err
	}

	return &models.TokenModel{
		ID:            user.ID.Hex(),
		Name:          user.FirstName + " " + user.LastName,
		Username:      user.Username,
		Email:         user.Email,
		Role:          user.Role,
		EmailVerified: user.EmailVerified,
		Active:        user.Active,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		Token: map[string]interface{}{
			"token":                 newAccessToken,
			"type":                  "Bearer",
			"expiresIn":             utils.TokenExpiryDuration,
			"expiresAt":             time.Now().Add(utils.TokenExpiryDuration),
			"refreshToken":          newRefreshToken,
			"refreshTokenExpiresIn": utils.RefreshTokenExpiryDuration,
			"refreshTokenExpiresAt": time.Now().Add(utils.RefreshTokenExpiryDuration),
		},
	}, nil
}
