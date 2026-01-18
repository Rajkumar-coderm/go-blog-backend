package blogs

import (
	"context"
	"time"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SavePost(c *gin.Context) {
	postsColl := config.DB.Collection("posts")
	savedColl := config.DB.Collection("bookmarked")

	// Get userID from context
	userID, ok := c.Get("userID")
	if !ok {
		utils.SendError(c, utils.ErrUnauthorized)
		return
	}

	// Parse request body
	var req struct {
		PostID string `json:"id" binding:"required"`
		Save   bool   `json:"save"` // true = save, false = unsave
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, utils.ErrBadRequest)
		return
	}

	userObjID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		utils.SendError(c, utils.ErrUnauthorized)
		return
	}

	postObjID, err := primitive.ObjectIDFromHex(req.PostID)
	if err != nil {
		utils.SendError(c, utils.ErrBadRequest)
		return
	}

	// Check if post exists
	if err := postsColl.FindOne(context.TODO(), bson.M{"_id": postObjID}).Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			utils.SendError(c, utils.ErrNotFound)
			return
		}
		utils.SendError(c, utils.ErrInternalServer)
		return
	}

	filter := bson.M{
		"userId": userObjID,
		"postId": postObjID,
	}

	// ---------------- SAVE ----------------
	if req.Save {
		count, err := savedColl.CountDocuments(context.TODO(), filter)
		if err != nil {
			utils.SendError(c, utils.ErrInternalServer)
			return
		}

		if count > 0 {
			utils.SendError(
				c,
				utils.NewAPIError(409, "This post is already saved"),
			)
			return
		}

		now := time.Now()
		doc := models.SavedPost{
			ID:        primitive.NewObjectID(),
			UserID:    userObjID,
			PostID:    postObjID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if _, err := savedColl.InsertOne(context.TODO(), doc); err != nil {
			utils.SendError(c, utils.ErrInternalServer)
			return
		}

		utils.SendSuccess(c, "Post saved successfully", nil, 201)
		return
	}

	// ---------------- UNSAVE ----------------
	res, err := savedColl.DeleteOne(context.TODO(), filter)
	if err != nil {
		utils.SendError(c, utils.ErrInternalServer)
		return
	}

	if res.DeletedCount == 0 {
		utils.SendError(
			c,
			utils.NewAPIError(400, "Post is not saved yet"),
		)
		return
	}

	c.Status(204)
}
