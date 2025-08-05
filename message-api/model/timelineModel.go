package model

import (
	"time"
)

type TimelineItem struct {
	MessageID string    `json:"message_id" dynamodbav:"message_id"`
	UserID    string    `json:"user_id" dynamodbav:"user_id"`
	AuthorID  string    `json:"author_id" dynamodbav:"author_id"`
	Content   string    `json:"content" dynamodbav:"content"`
	CreatedAt time.Time `json:"created_at" dynamodbav:"created_at"`
}
