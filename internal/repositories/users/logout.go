package users

import (
	"time"

	"github.com/Rajkumar-coderm/go-blog-backend/internal/auth"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/repositories/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func LogoutUser(c *gin.Context, userID primitive.ObjectID, refreshToken string, logoutAll bool) error {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		tokenStr := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenStr = authHeader[7:]
		}
		if tokenStr != "" {
			claims, err := auth.ValidateJWT(tokenStr)
			if err == nil {
				var exp time.Time
				if claims.ExpiresAt != nil {
					exp = claims.ExpiresAt.Time
				} else {
					exp = time.Now().Add(15 * time.Minute)
				}
				_ = sessions.AddTokenToBlacklist(tokenStr, exp)
			}
		}
	}

	if logoutAll {
		err := sessions.DeleteUserSessions(userID)
		if err != nil {
			return err
		}
	} else {
		err := sessions.DeleteSessionByToken(refreshToken)
		if err != nil {
			return err
		}
	}

	return nil
}
