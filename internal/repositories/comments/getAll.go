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
	// userCol := config.DB.Collection("users")

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

	var exitsPost models.Post

	_err := col.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&exitsPost)

	if _err != nil {
		return nil, 0, _err
	}

	cursor, err := col.Find(context.TODO(), bson.M{"_id": objID, "comments": bson.M{"$exists": true}})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	var result struct {
		Comments          []models.Comment `bson:"comments"`
		TotalCommentCount int64            `bson:"totalCommentCount"`
	}

	if cursor.Next(context.TODO()) {
		if err := cursor.Decode(&result); err != nil {
			return nil, 0, err
		}
	}

	return result.Comments, result.TotalCommentCount, nil
}
