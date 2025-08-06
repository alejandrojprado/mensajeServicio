package service

import (
	"context"
	"time"

	"mensajesService/components/database"
	"mensajesService/message-api/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type MessageServiceInterface interface {
	CreateMessage(ctx context.Context, userID, content string) (*model.Message, error)
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

	return message, nil
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

func generateUUID() string {
	return uuid.New().String()
}
