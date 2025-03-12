package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents the user schema in MongoDB
type User struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName           string             `bson:"first_name,omitempty" json:"firstName,omitempty"`
	LastName            string             `bson:"last_name,omitempty" json:"lastName,omitempty"`
	Email               string             `bson:"email" json:"email" validate:"required,email,unique"`
	Username            string             `bson:"username" json:"username" validate:"required,unique"`
	Password            string             `bson:"password" json:"password" validate:"required"`
	Role                string             `bson:"role" json:"role" validate:"oneof=user editor admin"`
	EmailVerified       bool               `bson:"email_verified" json:"emailVerified"`
	VerificationToken   string             `bson:"verification_token,omitempty" json:"verificationToken,omitempty"`
	ResetPasswordToken  string             `bson:"reset_password_token,omitempty" json:"resetPasswordToken,omitempty"`
	ResetPasswordExpiry *time.Time         `bson:"reset_password_expires,omitempty" json:"resetPasswordExpires,omitempty"`
	LastLogin           *time.Time         `bson:"last_login,omitempty" json:"lastLogin,omitempty"`
	Contact             UserContact        `bson:"contact,omitempty" json:"contact,omitempty"`
	Blog                UserBlog           `bson:"blog,omitempty" json:"blog,omitempty"`
	Activity            UserActivity       `bson:"activity,omitempty" json:"activity,omitempty"`
	Birthday            *time.Time         `bson:"birthday,omitempty" json:"birthday,omitempty"`
	Avatar              string             `bson:"avatar,omitempty" json:"avatar,omitempty"`
	Location            string             `bson:"location,omitempty" json:"location,omitempty"`
	Notifications       UserNotifications  `bson:"notifications,omitempty" json:"notifications,omitempty"`
	Bio                 string             `bson:"bio,omitempty" json:"bio,omitempty" validate:"max=500"`
	Token               string             `bson:"token,omitempty" json:"token,omitempty"`
	RefreshToken        string             `bson:"refresh_token,omitempty" json:"refreshToken,omitempty"`

	// System Fields
	Active    bool      `bson:"active" json:"active"`
	CreatedAt time.Time `bson:"created_at,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updatedAt,omitempty"`
}

// UserContact contains contact and social fields
type UserContact struct {
	Phone  string `bson:"phone,omitempty" json:"phone,omitempty"`
	Social struct {
		Twitter   string `bson:"twitter,omitempty" json:"twitter,omitempty"`
		Facebook  string `bson:"facebook,omitempty" json:"facebook,omitempty"`
		Instagram string `bson:"instagram,omitempty" json:"instagram,omitempty"`
		LinkedIn  string `bson:"linkedin,omitempty" json:"linkedin,omitempty"`
		GitHub    string `bson:"github,omitempty" json:"github,omitempty"`
	} `bson:"social,omitempty" json:"social"`
}

// UserBlog contains blog-related fields
type UserBlog struct {
	Articles  []primitive.ObjectID `bson:"articles,omitempty" json:"articles,omitempty"`
	Bookmarks []primitive.ObjectID `bson:"bookmarks,omitempty" json:"bookmarks,omitempty"`
	Following []primitive.ObjectID `bson:"following,omitempty" json:"following,omitempty"`
	Followers []primitive.ObjectID `bson:"followers,omitempty" json:"followers,omitempty"`
	Interests []string             `bson:"interests,omitempty" json:"interests,omitempty"`
}

// UserActivity contains activity and metrics fields
type UserActivity struct {
	LastActive      *time.Time `bson:"last_active,omitempty" json:"lastActive,omitempty"`
	CommentCount    int        `bson:"comment_count" json:"commentCount"`
	ArticleCount    int        `bson:"article_count" json:"articleCount"`
	UpvotesReceived int        `bson:"upvotes_received" json:"upvotesReceived"`
}

// UserNotifications contains notification settings
type UserNotifications struct {
	Email struct {
		Newsletter bool `bson:"newsletter" json:"newsletter"`
		Comments   bool `bson:"comments" json:"comments"`
		Followers  bool `bson:"followers" json:"followers"`
	} `bson:"email" json:"email"`
	PushEnabled bool `bson:"push_enabled" json:"pushEnabled"`
}