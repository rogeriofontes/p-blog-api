package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PostCategory representa uma categoria de posts no MongoDB.
type PostCategory struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
