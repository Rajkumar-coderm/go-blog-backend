package sessions

import (
	"context"
	"time"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateSession creates a new user session in the database
func CreateSession(userID primitive.ObjectID, refreshToken, deviceInfo, ipAddress string, expiresAt time.Time) (*models.UserSession, error) {
	col := config.DB.Collection("user_sessions")

	session := &models.UserSession{
		ID:           primitive.NewObjectID(),
		UserID:       userID,
		RefreshToken: refreshToken,
		DeviceInfo:   deviceInfo,
		IPAddress:    ipAddress,
		ExpiresAt:    expiresAt,
		CreatedAt:    time.Now(),
	}

	_, err := col.InsertOne(context.TODO(), session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// GetSessionByToken retrieves a session by refresh token
func GetSessionByToken(refreshToken string) (*models.UserSession, error) {
	col := config.DB.Collection("user_sessions")

	var session models.UserSession
	err := col.FindOne(context.TODO(), map[string]interface{}{
		"refresh_token": refreshToken,
		"expires_at": map[string]interface{}{
			"$gt": time.Now(),
		},
	}).Decode(&session)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &session, nil
}

// DeleteSession deletes a session by ID
func DeleteSession(sessionID primitive.ObjectID) error {
	col := config.DB.Collection("user_sessions")

	_, err := col.DeleteOne(context.TODO(), map[string]interface{}{
		"_id": sessionID,
	})

	return err
}

// DeleteSessionByToken deletes a session by refresh token
func DeleteSessionByToken(refreshToken string) error {
	col := config.DB.Collection("user_sessions")

	_, err := col.DeleteOne(context.TODO(), map[string]interface{}{
		"refresh_token": refreshToken,
	})

	return err
}

// DeleteUserSessions deletes all sessions for a user
func DeleteUserSessions(userID primitive.ObjectID) error {
	col := config.DB.Collection("user_sessions")

	_, err := col.DeleteMany(context.TODO(), map[string]interface{}{
		"user_id": userID,
	})

	return err
}
