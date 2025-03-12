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

// / Like/Dislike any post by post id and user id user id we are taking from token...
// / [userID] this is user id from token
// / [postID] this is post id from url
// / [like] this is true or false from request body
func LikeDislikePost(c *gin.Context) error {
	collection := config.DB.Collection("posts")

	// Get userID from token (context)
	userID, ok := c.Get("userID")
	if !ok {
		return errors.New("invalid user ID")
	}

	// Parse request body
	var request struct {
		Like bool   `json:"like"`
		ID   string `json:"id" binding:"required"`
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

	// Handle Like
	if request.Like {
		// Check if user already liked the post
		existingLike := collection.FindOne(context.TODO(), bson.M{
			"_id":   postObjID,
			"likes": bson.M{"$elemMatch": bson.M{"userId": userObjID}},
		})

		if existingLike.Err() == nil { // User already liked, update timestamp
			_, err := collection.UpdateOne(context.TODO(), bson.M{
				"_id":          postObjID,
				"likes.userId": userObjID,
			}, bson.M{
				"$set": bson.M{"likes.$.updatedAt": time.Now()},
			})
			return err
		} else if existingLike.Err() != mongo.ErrNoDocuments {
			return existingLike.Err()
		}

		// Add a new like entry
		likeEntry := models.UserAtraction{
			ID:        primitive.NewObjectID(),
			UserId:    userObjID,
			PostId:    postObjID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": postObjID}, bson.M{
			"$addToSet": bson.M{"likes": likeEntry},
		})
		return err
	}

	// Handle Dislike (Unlike)
	res, err := collection.UpdateOne(context.TODO(), bson.M{"_id": postObjID}, bson.M{
		"$pull": bson.M{"likes": bson.M{"userId": userObjID}},
	})
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("like entry not found")
	}

	return nil
}
