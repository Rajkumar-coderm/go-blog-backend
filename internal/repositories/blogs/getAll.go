package blogs

import (
	"context"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Get all posts with author details
func GetAll(ctx *gin.Context, request *models.GetRequest, userID string) (interface{}, int64, error) {
	col := config.DB.Collection("posts")

	query := bson.M{}
	if request.Q != "" {
		query["$text"] = bson.M{"$search": request.Q} // Enables text search
	}

	// Convert userID string to ObjectID
	var userObjectID primitive.ObjectID
	if objID, err := primitive.ObjectIDFromHex(userID); err == nil {
		userObjectID = objID
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: query}},

		// Lookup user details by userId
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "userId"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "auther"},
		}}},

		// Unwind author to convert array into an object
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$auther"},
			{Key: "preserveNullAndEmptyArrays", Value: true}, // Keeps posts with no user attached
		}}},

		bson.D{{Key: "$addFields", Value: bson.D{
			// Ensure `totalLikes` and `totalCommentCount` always work
			{Key: "totalLikes", Value: bson.D{
				{Key: "$size", Value: bson.D{
					{Key: "$ifNull", Value: bson.A{"$likes", bson.A{}}}, // Avoid null error
				}},
			}},
			{Key: "totalCommentCount", Value: bson.D{
				{Key: "$size", Value: bson.D{
					{Key: "$ifNull", Value: bson.A{"$comments", bson.A{}}}, // Avoid null error
				}},
			}},
			{Key: "auther.fullName", Value: bson.D{
				{Key: "$concat", Value: bson.A{"$auther.first_name", " ", "$auther.last_name"}},
			}},

			{Key: "auther.username", Value: bson.D{
				{Key: "$concat", Value: bson.A{"$auther.first_name", " ", "$auther.last_name"}},
			}},

			// `isLiked` to correctly check if the user liked the post
			{Key: "isLiked", Value: bson.D{
				{Key: "$gt", Value: bson.A{
					bson.D{{Key: "$size", Value: bson.D{{Key: "$filter", Value: bson.D{
						{Key: "input", Value: bson.D{{Key: "$ifNull", Value: bson.A{"$likes", bson.A{}}}}},
						{Key: "as", Value: "like"},
						{Key: "cond", Value: bson.D{{Key: "$eq", Value: bson.A{"$$like.userId", userObjectID}}}},
					}}}}},
					0,
				}},
			}},

			// `isBookmarked` to correctly check if the user bookmarked the post
			{Key: "isBookmarked", Value: bson.D{
				{Key: "$gt", Value: bson.A{
					bson.D{{Key: "$size", Value: bson.D{{Key: "$filter", Value: bson.D{
						{Key: "input", Value: bson.D{{Key: "$ifNull", Value: bson.A{"$bookmarks", bson.A{}}}}},
						{Key: "as", Value: "bookmark"},
						{Key: "cond", Value: bson.D{{Key: "$eq", Value: bson.A{"$$bookmark.userId", userObjectID}}}},
					}}}}},
					0,
				}},
			}},
		}}},

		// Project only required fields
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "title", Value: 1},
			{Key: "content", Value: 1},
			{Key: "totalLikes", Value: 1},
			{Key: "totalCommentCount", Value: 1},
			{Key: "created_at", Value: 1},
			{Key: "updated_at", Value: 1},
			{Key: "isLiked", Value: 1},
			{Key: "isBookmarked", Value: 1},
			{Key: "auther", Value: bson.D{
				{Key: "_id", Value: "$auther._id"},
				{Key: "username", Value: "$auther.username"},
				{Key: "fullName", Value: "$auther.fullName"},
				{Key: "bio", Value: "$auther.bio"},
			}},
		}}},

		// Pagination
		bson.D{{Key: "$skip", Value: request.Skip}},
		bson.D{{Key: "$limit", Value: request.Limit}},
	}

	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	var posts []bson.M
	if err := cursor.All(context.TODO(), &posts); err != nil {
		return nil, 0, err
	}

	totalCount, countErr := col.CountDocuments(context.TODO(), query)
	if countErr != nil {
		return nil, 0, countErr
	}

	return posts, totalCount, nil
}

// Get a post by ID
func GetPostByID(c *gin.Context) (*models.Post, error) {
	col := config.DB.Collection("posts")

	objID, err := primitive.ObjectIDFromHex(c.Query("id"))
	if err != nil {
		return nil, err
	}

	var post models.Post
	err = col.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&post)
	if err != nil {
		return nil, err
	}
	post.TotalLikes = len(post.Like)
	return &post, nil
}
