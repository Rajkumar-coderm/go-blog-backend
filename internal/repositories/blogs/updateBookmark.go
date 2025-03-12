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
	collection := config.DB.Collection("posts")

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
	err = collection.FindOne(context.TODO(), bson.M{"_id": postObjID}).Err()
	if err == mongo.ErrNoDocuments {
		return errors.New("post not found")
	} else if err != nil {
		return err
	}

	if request.IsBookmark {
		// Check if user already bookmarked the post
		var post struct {
			Bookmarks []models.UserAtraction `bson:"bookmarks"`
		}
		err := collection.FindOne(context.TODO(), bson.M{"_id": postObjID}).Decode(&post)
		if err != nil && err != mongo.ErrNoDocuments {
			return err
		}

		// Check if user has already bookmarked
		for _, bookmark := range post.Bookmarks {
			if bookmark.UserId == userObjID {
				// Update existing bookmark timestamp
				_, err := collection.UpdateOne(context.TODO(), bson.M{
					"_id":              postObjID,
					"bookmarks.userId": userObjID,
				}, bson.M{
					"$set": bson.M{"bookmarks.$.updatedAt": time.Now()},
				})
				return err
			}
		}

		// Add a new bookmark entry
		bookmarkEntry := models.UserAtraction{
			ID:        primitive.NewObjectID(),
			UserId:    userObjID,
			PostId:    postObjID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": postObjID}, bson.M{
			"$push": bson.M{"bookmarks": bookmarkEntry},
		})
		return err
	}

	// Handle Remove Bookmark
	res, err := collection.UpdateOne(context.TODO(), bson.M{"_id": postObjID}, bson.M{
		"$pull": bson.M{"bookmarks": bson.M{"userId": userObjID}},
	})

	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("bookmark entry not found")
	}

	return nil
}
