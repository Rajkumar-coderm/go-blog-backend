package models

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;index" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `gorm:"index" json:"deletedAt"`
} //@name Model

type APIResponse struct {
	Status     string      `json:"status" example:"success"`
	Message    string      `json:"message,omitempty" example:"success"`
	ErrorCode  string      `json:"errorCode,omitempty" example:"0"`
	Data       interface{} `json:"data,omitempty"`
	TotalCount int64       `json:"totalCount,omitempty" example:"10"`
} //@name APIResponse

// ErrorResponse is the common error response body to be used in case of any error
type ErrorResponse struct {
	GroupId  string `json:"groupId" example:"0"`
	ClientId string `json:"clientId" example:"5"`
	UserId   string `json:"userId" example:"5"`
} //@name ErrorResponse

// Data from API Gateway
type UserDataFromAPIGateWay struct {
	ClientId         string      `json:"clientId" example:"5"`
	ClientName       string      `json:"clientName" example:"5"`
	UserId           string      `json:"userId" example:"5"`
	Token            string      `json:"token" example:"5"`
	RealmNameOfToken string      `json:"realmNameOfToken" example:"master"`
	AuthCompleted    bool        `json:"authCompleted" example:"true"`
	MetaData         interface{} `json:"metaData" example:"{}"`
} //@name UserDataFromAPIGateWay

// {"clientId": "", "userId": "", "token":"", "authCompleted":true, "realmNameOfToken":"master"}

type IdRequest struct {
	ID int `json:"id" binding:"required" example:"1"`
} //@name IdRequest

type DeleteRequest struct {
	Ids []string `json:"ids" binding:"required" example:"[1,2,3]"`
} //@name DeleteRequest

type PatchStatusRequest struct {
	ID     primitive.ObjectID `json:"id" bson:"_id" binding:"required" example:"5f5f5f5f5f5f5f5f5f5f5f5f"`
	Status string             `json:"status" bson:"status" binding:"required" example:"active"` //User status Active, Deleted, Suspended
	Reason string             `json:"reason" bson:"reason" example:"not verified"`
} //@name PatchStatusRequest

type GetRequest struct {
	Skip   int    `json:"skip" example:"0" form:"skip"`    // records to skip
	Limit  int    `json:"limit" example:"10" form:"limit"` // limit for the records. non of records to fetch per page
	Id     string `json:"id" form:"id"`
	Status string `json:"status" form:"status"`
	Q      string `json:"q" form:"q"`
} //@name GetRequest

type Stamps struct {
	CreatedAt int64 `json:"createdAt" bson:"createdAt" example:"1600000000"`
	UpdatedAt int64 `json:"updatedAt" bson:"updatedAt" example:"1600000000"`
} //@name Stamps

type GetVerificationDataReq struct {
	Token string `json:"token" form:"token" binding:"required"`
} //@name GetVerificationDataReq

type Address struct {
	Address      string `json:"address" bson:"address" binding:"required"`
	AddressLine2 string `json:"addressLine2" bson:"addressLine2"`
	City         string `json:"city" bson:"city" binding:"required"`
	State        string `json:"state" bson:"state" binding:"required"`
	Country      string `json:"country" bson:"country" binding:"required"`
	ZipCode      string `json:"zipCode" bson:"zipCode" binding:"required"`
} // @name Address

type StatusLog struct {
	Status    string `json:"status" bson:"status"`
	TimeStamp int64  `json:"timeStamp" bson:"timeStamp"`
	UserRole  string `json:"userRole" bson:"userRole"`
	UserId    string `json:"userId" bson:"userId"`
	UserName  string `json:"userName" bson:"userName"`
	Reason    string `json:"reason,omitempty" bson:"reason,omitempty"`
} // @name StatusLog

type Count struct {
	TotalCount int `json:"totalCount" bson:"totalCount"`
} // @name Count
