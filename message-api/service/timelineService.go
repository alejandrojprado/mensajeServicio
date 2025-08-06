package service

import (
	"context"

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
	UpdateFollowersTimeline(ctx context.Context, message *model.Message) error
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

func (s *TimelineService) UpdateFollowersTimeline(ctx context.Context, message *model.Message) error {
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

func (s *TimelineService) getFollowers(ctx context.Context, userID string) ([]string, error) {
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

func (s *TimelineService) saveTimelineItem(ctx context.Context, item *model.TimelineItem) error {
	timelineItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	logger.LogInfo("Saving timeline item", "message_id", item.MessageID, "user_id", item.UserID, "author_id", item.AuthorID)

	timelineItem["PK"] = &types.AttributeValueMemberS{Value: "USER#" + item.UserID}
	timelineItem["SK"] = &types.AttributeValueMemberS{Value: "MESSAGE#" + item.MessageID}

	logger.LogInfo("Timeline item keys", "pk", "USER#"+item.UserID, "sk", "MESSAGE#"+item.MessageID)

	return s.dbClient.PutItem(ctx, s.dbClient.GetTimelineTableName(), timelineItem)
}
