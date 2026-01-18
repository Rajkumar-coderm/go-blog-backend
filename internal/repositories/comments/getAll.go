package comments

import (
	"context"
	"errors"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllPostComments(ctx *gin.Context, request *models.GetRequest) (interface{}, int64, error) {
	col := config.DB.Collection("posts")
	commentsCol := config.DB.Collection("comments")

	postID := ctx.Query("id")
	if postID == "" {
		return nil, 0, errors.New("id is required")
	}

	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, 0, err
	}

	if request.Skip < 0 {
		request.Skip = 0
	}
	if request.Limit < 1 {
		request.Limit = 40
	}

	var existsPost models.Post

	_err := col.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&existsPost)

	if _err != nil {
		return nil, 0, _err
	}

	// Fetch all comments for the post
	cursor, err := commentsCol.Find(context.TODO(), bson.M{"postId": objID})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	var allComments []models.Comment
	if err = cursor.All(context.TODO(), &allComments); err != nil {
		return nil, 0, err
	}

	// Build the comment tree
	commentMap := make(map[primitive.ObjectID]*models.Comment)
	var rootComments []models.Comment

	for i := range allComments {
		comment := &allComments[i]
		commentMap[comment.ID] = comment
		if comment.ParentId == nil {
			rootComments = append(rootComments, *comment)
		} else {
			if parent, exists := commentMap[*comment.ParentId]; exists {
				parent.Replies = append(parent.Replies, comment)
			}
		}
	}

	return rootComments, int64(len(allComments)), nil
}
