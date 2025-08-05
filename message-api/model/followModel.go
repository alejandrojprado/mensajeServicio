package model

import (
	"time"
)

type Follow struct {
	FollowerID  string    `json:"follower_id" dynamodbav:"follower_id"`
	FollowingID string    `json:"following_id" dynamodbav:"following_id"`
	CreatedAt   time.Time `json:"created_at" dynamodbav:"created_at"`
}

type FollowRequest struct {
	FollowingID string `json:"following_id"`
}
