package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FavoritePost struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID  primitive.ObjectID `json:"user_id" bson:"user_id"`
	PostID  primitive.ObjectID `json:"post_id" bson:"post_id"`
	SavedAt time.Time          `json:"saved_at" bson:"saved_at"`
}
