package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/Rajkumar-coderm/go-blog-backend/config"
	"github.com/Rajkumar-coderm/go-blog-backend/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func VerifyUserName(c *gin.Context) (*models.User, error) {
	col := config.DB.Collection("users")
	var user models.User
	err := col.FindOne(context.TODO(), bson.M{"username": c.Query("username")}).Decode(&user)
	fmt.Println(err == nil)
	if err != nil {
		return &user, nil
	}
	return nil, errors.New("this username already exist please try another username")
}
