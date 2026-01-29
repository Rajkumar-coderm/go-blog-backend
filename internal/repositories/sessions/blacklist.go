package sessions

import (
	"context"
	"time"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddTokenToBlacklist stores an access token to the blacklist until its expiry
func AddTokenToBlacklist(token string, expiresAt time.Time) error {
	col := config.DB.Collection("blacklisted_tokens")

	entry := &models.BlacklistedToken{
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	_, err := col.InsertOne(context.TODO(), entry)
	return err
}

// IsTokenBlacklisted checks whether the given token exists in the blacklist and is not expired
func IsTokenBlacklisted(token string) (bool, error) {
	col := config.DB.Collection("blacklisted_tokens")

	var entry models.BlacklistedToken
	err := col.FindOne(context.TODO(), bson.M{"token": token, "expires_at": bson.M{"$gt": time.Now()}}).Decode(&entry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
