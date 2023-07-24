package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateReviewRequest struct {
	Name       string             `json:"name" bson:"name" binding:"required"`
	Rating     string             `json:"rating" bson:"rating" binding:"required"`
	Comment    string             `json:"comment,omitempty" bson:"comment,omitempty" binding:"required"`
	UserID     primitive.ObjectID `bson:"user_id"`
	ProductID  string             `json:"product_id,omitempty" bson:"product_id,omitempty"`
	Timestamps bool               `json:"timestamps,omitempty" bson:"timestamps,omitempty"`
	CreateAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type DBReview struct {
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	Rating     string             `json:"rating,omitempty" bson:"rating,omitempty" `
	Comment    string             `json:"comment,omitempty" bson:"comment,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id"`
	ProductID  string             `json:"product_id,omitempty" bson:"product_id,omitempty"`
	Timestamps bool               `json:"timestamps,omitempty" bson:"timestamps,omitempty"`
	CreateAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UpdateReview struct {
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	Rating     string             `json:"rating,omitempty" bson:"rating,omitempty" `
	Comment    string             `json:"comment,omitempty" bson:"comment,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id"`
	ProductID  string             `json:"product_id,omitempty" bson:"product_id,omitempty"`
	Timestamps bool               `json:"timestamps,omitempty" bson:"timestamps,omitempty"`
	CreateAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
