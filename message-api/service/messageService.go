package service

import (
	"context"
	"time"

	"mensajesService/components/database"
	"mensajesService/components/logger"
	"mensajesService/message-api/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type MessageServiceInterface interface {
	CreateMessage(ctx context.Context, userID, content string) (*model.Message, error)
	GetUserTimeline(ctx context.Context, userID string, limit int) ([]*model.TimelineItem, error)
	GetUserMessages(ctx context.Context, userID string, limit int) ([]*model.Message, error)
}

type MessageService struct {
	dbClient database.DDBClientInterface
}

func NewMessageService(dbClient database.DDBClientInterface) *MessageService {
	return &MessageService{dbClient: dbClient}
}

func (s *MessageService) CreateMessage(ctx context.Context, userID, content string) (*model.Message, error) {
	messageID := generateUUID()
	now := time.Now()

	message := &model.Message{
		ID:        messageID,
		UserID:    userID,
		Content:   content,
		CreatedAt: now,
	}

	err := s.saveMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := s.updateFollowersTimeline(context.Background(), message); err != nil {
			logger.LogError("Error updating followers timeline", "error", err, "message_id", messageID)
		}
	}()

	return message, nil
}

func (s *MessageService) GetUserTimeline(ctx context.Context, userID string, limit int) ([]*model.TimelineItem, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(s.dbClient.GetTimelineTableName()),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "USER#" + userID},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(int32(limit)),
	}

	result, err := s.dbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	var timelineItems []*model.TimelineItem
	err = attributevalue.UnmarshalListOfMaps(result.Items, &timelineItems)
	if err != nil {
		return nil, err
	}

	return timelineItems, nil
}

func (s *MessageService) GetUserMessages(ctx context.Context, userID string, limit int) ([]*model.Message, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(s.dbClient.GetMessagesTableName()),
		IndexName:              aws.String("UserIDIndex"),
		KeyConditionExpression: aws.String("user_id = :user_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_id": &types.AttributeValueMemberS{Value: userID},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(int32(limit)),
	}

	result, err := s.dbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	var messages []*model.Message
	err = attributevalue.UnmarshalListOfMaps(result.Items, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *MessageService) saveMessage(ctx context.Context, message *model.Message) error {
	item, err := attributevalue.MarshalMap(message)
	if err != nil {
		return err
	}

	item["PK"] = &types.AttributeValueMemberS{Value: "MESSAGE#" + message.ID}
	item["SK"] = &types.AttributeValueMemberS{Value: "MESSAGE#" + message.ID}

	return s.dbClient.PutItem(ctx, s.dbClient.GetMessagesTableName(), item)
}

func (s *MessageService) updateFollowersTimeline(ctx context.Context, message *model.Message) error {
	followers, err := s.getFollowers(ctx, message.UserID)
	if err != nil {
		return err
	}

	for _, followerID := range followers {
		timelineItem := &model.TimelineItem{
			MessageID: message.ID,
			UserID:    followerID,
			AuthorID:  message.UserID,
			Content:   message.Content,
			CreatedAt: message.CreatedAt,
		}

		err := s.saveTimelineItem(ctx, timelineItem)
		if err != nil {
			logger.LogError("Error saving timeline item", "error", err, "follower_id", followerID)
			continue
		}
	}

	return nil
}

func (s *MessageService) getFollowers(ctx context.Context, userID string) ([]string, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(s.dbClient.GetFollowersTableName()),
		IndexName:              aws.String("FollowingIndex"),
		KeyConditionExpression: aws.String("SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sk": &types.AttributeValueMemberS{Value: "FOLLOWING#" + userID},
		},
	}

	result, err := s.dbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	var follows []*model.Follow
	err = attributevalue.UnmarshalListOfMaps(result.Items, &follows)
	if err != nil {
		return nil, err
	}

	var followers []string
	for _, follow := range follows {
		followers = append(followers, follow.FollowerID)
	}

	return followers, nil
}

func (s *MessageService) saveTimelineItem(ctx context.Context, item *model.TimelineItem) error {
	timelineItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}
	timelineItem["PK"] = &types.AttributeValueMemberS{Value: "USER#" + item.UserID}
	timelineItem["SK"] = &types.AttributeValueMemberS{Value: "MESSAGE#" + item.MessageID}

	return s.dbClient.PutItem(ctx, s.dbClient.GetTimelineTableName(), timelineItem)
}

func generateUUID() string {
	return uuid.New().String()
}
