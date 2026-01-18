package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post represents a blog post structure
type Post struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title          string             `bson:"title" json:"title"`
	Content        string             `bson:"content" json:"content"`
	UserId         primitive.ObjectID `bson:"userId" json:"author"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at,omitempty" json:"updated_at"`
	CommentsCount  int                `bson:"commentsCount" json:"commentsCount"`
	TotalLikes     int                `bson:"totalLikes" json:"totalLikes"`
	Like           []UserAtraction    `bson:"likes,omitempty" json:"likes,omitempty"`
	BookmarksCount int                `bson:"bookmarksCount" json:"bookmarksCount"`
	IsLiked        bool               `bson:"isLiked,omitempty" json:"isLiked,omitempty"`
	IsBookmarked   bool               `bson:"isBookmarked,omitempty" json:"isBookmarked,omitempty"`
}

type Comment struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserId    primitive.ObjectID  `bson:"userId" json:"userId"`
	PostId    primitive.ObjectID  `bson:"postId" json:"postId"`
	ParentId  *primitive.ObjectID `bson:"parentId,omitempty" json:"parentId,omitempty"`
	Content   string              `bson:"content" json:"content"`
	CreatedAt time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time           `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	Replies   []*Comment          `bson:"replies,omitempty" json:"replies,omitempty"`
	Auther    Auther              `bson:"auther,omitempty" json:"auther,omitempty"`
}

type UserAtraction struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId    primitive.ObjectID `bson:"userId" json:"userId"`
	PostId    primitive.ObjectID `bson:"postId" json:"postId"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type Auther struct {
	UserId   primitive.ObjectID `bson:"userId" json:"userId"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	UserName string             `bson:"userName" json:"userName"`
	Bio      string             `bson:"bio" json:"bio"`
}
