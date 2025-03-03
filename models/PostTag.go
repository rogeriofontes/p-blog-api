package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PostTag representa uma tag no MongoDB.
type PostTag struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
