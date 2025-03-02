package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title      string             `json:"title" bson:"title"`
	CategoryID primitive.ObjectID `json:"category_id" bson:"category_id" binding:"required"` // Agora é ObjectID
	Category   *PostCategory      `json:"category,omitempty" bson:"category,omitempty"`
	Content    string             `json:"content" bson:"content"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}
