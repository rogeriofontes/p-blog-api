package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Follower struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	FollowID   primitive.ObjectID `json:"follow_id" bson:"follow_id"`
	FollowedAt time.Time          `json:"followed_at" bson:"followed_at"`
}
