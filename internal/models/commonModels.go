package models

import "time"

type CommonGetResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// User represents a user in the system.
type TokenModel struct {
	ID            string      `json:"id" bson:"_id"`
	Name          string      `json:"name" bson:"name"`
	Username      string      `json:"username" bson:"username"`
	Email         string      `json:"email" bson:"email"`
	Role          string      `json:"role" bson:"role"`
	EmailVerified bool        `json:"emailVerified" bson:"emailVerified"`
	Contact       string      `json:"contact" bson:"contact"`
	Activity      string      `json:"activity" bson:"activity"`
	Token         interface{} `json:"token" bson:"token"`
	Active        bool        `json:"active" bson:"active"`
	CreatedAt     time.Time   `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt" bson:"updatedAt"`
}
