package service

import (
	"context"

	"mensajesService/components/database"
	"mensajesService/message-api/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TimelineServiceInterface interface {
	GetUserTimeline(ctx context.Context, userID string, limit int) ([]*model.TimelineItem, error)
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
