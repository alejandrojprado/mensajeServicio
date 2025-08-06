package model

import (
	"time"
)

type Message struct {
	ID        string    `json:"id" dynamodbav:"message_id"`
	UserID    string    `json:"user_id" dynamodbav:"user_id"`
	Content   string    `json:"content" dynamodbav:"content"`
	CreatedAt time.Time `json:"created_at" dynamodbav:"created_at"`
}
