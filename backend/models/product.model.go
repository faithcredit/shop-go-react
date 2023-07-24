package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateProductRequest struct {
	Name        string    `json:"name" bson:"name" binding:"required"`
	Price       string    `json:"price" bson:"price" binding:"required"`
	Category    string    `json:"category,omitempty" bson:"category,omitempty"`
	Brand       string    `json:"brand" bson:"brand" binding:"required"`
	Description string    `json:"description" bson:"description" binding:"required"`
	Image       string    `json:"image,omitempty" bson:"image,omitempty"`
	Qty         string    `json:"qty,omitempty" bson:"qty,omitempty"`
	CreateAt    time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type DBProduct struct {
	Id          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Price       string             `json:"price,omitempty" bson:"price,omitempty"`
	Category    string             `json:"category,omitempty" bson:"category,omitempty"`
	Brand       string             `json:"brand,omitempty" bson:"brand,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Image       string             `json:"image,omitempty" bson:"image,omitempty"`
	Qty         string             `json:"qty,omitempty" bson:"qty,omitempty"`
	CreateAt    time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UpdateProduct struct {
	Id          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Price       string             `json:"price,omitempty" bson:"price,omitempty"`
	Category    string             `json:"category,omitempty" bson:"category,omitempty"`
	Brand       string             `json:"brand,omitempty" bson:"brand,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Image       string             `json:"image,omitempty" bson:"image,omitempty"`
	Qty         string             `json:"qty,omitempty" bson:"qty,omitempty"`
	CreateAt    time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
