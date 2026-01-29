package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents the user schema in MongoDB
type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName     string             `bson:"first_name,omitempty" json:"firstName,omitempty"`
	LastName      string             `bson:"last_name,omitempty" json:"lastName,omitempty"`
	Email         string             `bson:"email" json:"email"`
	Username      string             `bson:"username" json:"username"`
	Password      string             `bson:"password" json:"password"`
	Role          string             `bson:"role" json:"role" validate:"oneof=user editor admin"`
	EmailVerified bool               `bson:"email_verified" json:"emailVerified"`
	LastLogin     *time.Time         `bson:"last_login,omitempty" json:"lastLogin,omitempty"`
	Birthday      *time.Time         `bson:"birthday,omitempty" json:"birthday,omitempty"`
	Avatar        string             `bson:"avatar,omitempty" json:"avatar,omitempty"`
	Location      string             `bson:"location,omitempty" json:"location,omitempty"`

	Bio              string    `bson:"bio,omitempty" json:"bio,omitempty" validate:"max=500"`
	Active           bool      `bson:"active" json:"active"`
	CreatedAt        time.Time `bson:"created_at,omitempty" json:"createdAt,omitempty"`
	UpdatedAt        time.Time `bson:"updated_at,omitempty" json:"updatedAt,omitempty"`
	RegistrationType string    `bson:"registration_type" json:"registrationType" validate:"oneof=email phone google"`
	PhoneVerified    bool      `bson:"phone_verified" json:"phoneVerified"`
	GoogleID         string    `bson:"google_id,omitempty" json:"googleId,omitempty"`
	Phone            string    `bson:"phone,omitempty" json:"phone,omitempty"`
	CoutryCode       string    `bson:"country_code,omitempty" json:"countryCode,omitempty"`
	PhoneIsoCode     string    `bson:"phone_iso_code,omitempty" json:"phoneIsoCode,omitempty"`
}

type LoginRequest struct {
	LoginType    string `json:"loginType" binding:"required,oneof=email phone google"`
	Email        string `json:"email" binding:"required_if=LoginType email,omitempty"`
	Password     string `json:"password" binding:"required"`
	Phone        string `json:"phone" binding:"required_if=LoginType phone,omitempty"`
	CountryCode  string `json:"countryCode" binding:"required_if=LoginType phone,omitempty"`
	PhoneIsoCode string `json:"phoneIsoCode" binding:"required_if=LoginType phone,omitempty"`
	GoogleID     string `json:"googleId" binding:"required_if=LoginType google,omitempty"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
	LogoutAll    bool   `json:"logoutAll"`
}
