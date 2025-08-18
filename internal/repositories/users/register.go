package users

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/auth"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Validation errors
var (
	ErrRegistrationTypeRequired = errors.New("registration type is required")
	ErrInvalidRegistrationType  = errors.New("invalid registration type")
	ErrEmailRequired            = errors.New("email is required for email registration")
	ErrPasswordRequired         = errors.New("password is required for email registration")
	ErrPhoneRequired            = errors.New("phone number is required for phone registration")
	ErrCountryCodeRequired      = errors.New("country code is required for phone registration")
	ErrPhoneIsoCodeRequired     = errors.New("phone iso code is required for phone registration")
	ErrGoogleIDRequired         = errors.New("google ID is required for google registration")
	ErrInvalidGoogleIDFormat    = errors.New("invalid google ID format")
	ErrInvalidGoogleIDChars     = errors.New("google ID contains invalid characters")
	ErrEmailExists              = errors.New("email already exists")
	ErrPhoneExists              = errors.New("phone number already exists")
	ErrGoogleAccountExists      = errors.New("google account already exists")
	ErrUsernameExists           = errors.New("username already exists")
)

// validateRegistrationData validates user registration data based on registration type
func validateRegistrationData(user *models.User) error {
	switch user.RegistrationType {
	case "email":
		if user.Email == "" {
			return ErrEmailRequired
		}
		if user.Password == "" {
			return ErrPasswordRequired
		}
	case "phone":
		if user.Phone == "" {
			return ErrPhoneRequired
		}
		if user.CoutryCode == "" {
			return ErrCountryCodeRequired
		}
		if user.PhoneIsoCode == "" {
			return ErrPhoneIsoCodeRequired
		}
		if user.Password == "" {
			return ErrPasswordRequired
		}
	case "google":
		if user.GoogleID == "" {
			return ErrGoogleIDRequired
		}
		// Validate Google ID format
		if len(user.GoogleID) < 21 || len(user.GoogleID) > 255 {
			return ErrInvalidGoogleIDFormat
		}
		// Validate Google ID characters
		for _, char := range user.GoogleID {
			if !unicode.IsLetter(char) && !unicode.IsNumber(char) && char != '.' {
				return ErrInvalidGoogleIDChars
			}
		}
		user.Password = ""
	default:
		return ErrInvalidRegistrationType
	}
	return nil
}

// checkExistingUser checks if a user already exists based on registration type
func checkExistingUser(col *mongo.Collection, user *models.User) error {
	var filter bson.M

	switch user.RegistrationType {
	case "email":
		filter = bson.M{"email": user.Email}
	case "phone":
		filter = bson.M{"phone": user.Phone}
	case "google":
		filter = bson.M{"google_id": user.GoogleID}
	}

	var existingUser models.User
	err := col.FindOne(context.TODO(), filter).Decode(&existingUser)
	if err == nil {
		switch user.RegistrationType {
		case "email":
			return ErrEmailExists
		case "phone":
			return ErrPhoneExists
		case "google":
			return ErrGoogleAccountExists
		}
	} else if err != mongo.ErrNoDocuments {
		return err
	}

	return nil
}

// generateUsername generates a username if not provided
func generateUsername(user *models.User) string {
	if user.Username != "" {
		return user.Username
	}

	switch user.RegistrationType {
	case "email":
		return strings.Split(user.Email, "@")[0]
	case "phone":
		return "user_" + user.Phone
	case "google":
		return "user_" + user.GoogleID
	default:
		return "user_" + user.ID.Hex()[:8]
	}
}

// createTokenResponse creates the token response model
func createTokenResponse(user *models.User, token, refreshToken string) *models.TokenModel {
	return &models.TokenModel{
		ID:            user.ID.Hex(),
		Name:          user.FirstName + " " + user.LastName,
		Username:      user.Username,
		Email:         user.Email,
		Role:          user.Role,
		EmailVerified: user.EmailVerified,
		Contact:       user.Contact.Phone,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		Token: map[string]interface{}{
			"token":                 token,
			"type":                  "Bearer",
			"expiresIn":             24 * time.Hour,
			"expiresAt":             time.Now().Add(24 * time.Hour),
			"refreshToken":          refreshToken,
			"refreshTokenExpiresIn": 7 * 24 * time.Hour,
			"refreshTokenExpiresAt": time.Now().Add(7 * 24 * time.Hour),
		},
	}
}

func RegisterUser(c *gin.Context, user *models.User) (*models.TokenModel, error) {
	col := config.DB.Collection("users")

	// Initialize user fields
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Validate registration type
	if user.RegistrationType == "" {
		return nil, ErrRegistrationTypeRequired
	}

	// Validate registration data
	if err := validateRegistrationData(user); err != nil {
		return nil, err
	}

	// Check if user already exists
	if err := checkExistingUser(col, user); err != nil {
		return nil, err
	}

	// Generate username if not provided
	user.Username = generateUsername(user)

	// Check if username already exists
	var existingUser models.User
	err := col.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		return nil, ErrUsernameExists
	} else if err != mongo.ErrNoDocuments {
		return nil, err
	}

	// Hash password if provided
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	// Insert user into database
	_, err = col.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}

	// Generate JWT auth token
	token, err := auth.GenerateJWT(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := auth.GenerateRefreshToken(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	// Create and return token response
	return createTokenResponse(user, token, refreshToken), nil
}
