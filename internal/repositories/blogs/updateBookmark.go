package blogs

import (
	"context"
	"errors"
	"time"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func BookmarkPost(c *gin.Context) error {
	postsColl := config.DB.Collection("posts")
	savedColl := config.DB.Collection("bookmarked")

	// Get userID from token (context)
	userID, ok := c.Get("userID")
	if !ok {
		return errors.New("invalid user ID")
	}

	// Parse request body
	var request struct {
		IsBookmark bool   `json:"isBookMark"`
		ID         string `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		return errors.New("invalid request payload")
	}

	// Convert IDs to ObjectID
	userObjID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		return errors.New("invalid user ID format")
	}

	postObjID, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		return errors.New("invalid post ID format")
	}

	// Check if the post exists
	err = postsColl.FindOne(context.TODO(), bson.M{"_id": postObjID}).Err()
	if err == mongo.ErrNoDocuments {
		return errors.New("post not found")
	} else if err != nil {
		return err
	}

	filter := bson.M{
		"userId": userObjID,
		"postId": postObjID,
	}

	if request.IsBookmark {
		// Check if already saved
		count, err := savedColl.CountDocuments(context.TODO(), filter)
		if err != nil {
			return err
		}
		if count > 0 {
			// Update timestamp
			_, err := savedColl.UpdateOne(context.TODO(), filter, bson.M{
				"$set": bson.M{"updatedAt": time.Now()},
			})
			return err
		}

		// Add new bookmark
		doc := models.SavedPost{
			ID:        primitive.NewObjectID(),
			UserID:    userObjID,
			PostID:    postObjID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if _, err := savedColl.InsertOne(context.TODO(), doc); err != nil {
			return err
		}

		// Increment bookmarksCount
		_, err = postsColl.UpdateOne(context.TODO(), bson.M{"_id": postObjID}, bson.M{
			"$inc": bson.M{"bookmarksCount": 1},
		})
		return err
	}

	// Remove bookmark
	res, err := savedColl.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("bookmark entry not found")
	}

	// Decrement bookmarksCount
	_, err = postsColl.UpdateOne(context.TODO(), bson.M{"_id": postObjID}, bson.M{
		"$inc": bson.M{"bookmarksCount": -1},
	})
	return err
}
