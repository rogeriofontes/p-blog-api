package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostReaction struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	PostID    primitive.ObjectID `json:"post_id" bson:"post_id"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Reaction  bool               `json:"reaction" bson:"reaction"` // "like" ou "dislike"
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
