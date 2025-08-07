package service

import (
	"context"
	"fmt"

	"mensajesService/components/database"
	"mensajesService/components/logger"
	"mensajesService/message-api/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TimelineServiceInterface interface {
	GetUserTimeline(ctx context.Context, userID string, limit int) ([]*model.TimelineItem, error)
	AddMessageToFollowersTimeline(ctx context.Context, message *model.Message, followerIDs []string) error
}

type TimelineService struct {
	dbClient database.DDBClientInterface
}

func NewTimelineService(dbClient database.DDBClientInterface) *TimelineService {
	return &TimelineService{dbClient: dbClient}
}

func (s *TimelineService) GetUserTimeline(ctx context.Context, userID string, limit int) ([]*model.TimelineItem, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(s.dbClient.GetTimelineTableName()),
		KeyConditionExpression: aws.String("user_id = :user_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_id": &types.AttributeValueMemberS{Value: userID},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(int32(limit)),
	}

	result, err := s.dbClient.Query(ctx, input)
	if err != nil {
		logger.LogError("Error getting user timeline", "error", err, "user_id", userID)
		return nil, err
	}

	var timelineItems []*model.TimelineItem
	err = attributevalue.UnmarshalListOfMaps(result.Items, &timelineItems)
	if err != nil {
		return nil, err
	}

	logger.LogInfo("Timeline retrieved successfully", "user_id", userID, "items_count", len(timelineItems))
	return timelineItems, nil
}

func (s *TimelineService) AddMessageToFollowersTimeline(ctx context.Context, message *model.Message, followerIDs []string) error {
	for _, followerID := range followerIDs {
		timelineItem := &model.TimelineItem{
			MessageID: message.ID,
			UserID:    followerID,
			AuthorID:  message.UserID,
			Content:   message.Content,
			CreatedAt: message.CreatedAt,
		}

		if err := s.saveTimelineItem(ctx, timelineItem); err != nil {
			logger.LogError("Error saving timeline item", "error", err, "message_id", message.ID, "follower_id", followerID)
			continue
		}
	}

	return nil
}

func (s *TimelineService) saveTimelineItem(ctx context.Context, item *model.TimelineItem) error {
	timelineItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	timelineItem["user_id"] = &types.AttributeValueMemberS{Value: item.UserID}
	timelineItem["timestamp"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", item.CreatedAt.Unix())}

	err = s.dbClient.PutItem(ctx, s.dbClient.GetTimelineTableName(), timelineItem)
	if err != nil {
		logger.LogError("Error saving timeline item", "error", err, "message_id", item.MessageID, "user_id", item.UserID)
		return err
	}

	return nil
}
