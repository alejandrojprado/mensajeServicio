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
)

type FollowServiceInterface interface {
	FollowUser(ctx context.Context, userID, followingID string) error
}

type FollowService struct {
	dbClient        database.DDBClientInterface
	messageService  MessageServiceInterface
	timelineService TimelineServiceInterface
}

func NewFollowService(dbClient database.DDBClientInterface, messageService MessageServiceInterface, timelineService TimelineServiceInterface) *FollowService {
	return &FollowService{
		dbClient:        dbClient,
		messageService:  messageService,
		timelineService: timelineService,
	}
}

func (s *FollowService) FollowUser(ctx context.Context, userID, followingID string) error {
	now := time.Now()

	follow := &model.Follow{
		FollowerID:  userID,
		FollowingID: followingID,
		CreatedAt:   now,
	}

	item, err := attributevalue.MarshalMap(follow)
	if err != nil {
		return err
	}

	item["follower_id"] = &types.AttributeValueMemberS{Value: userID}
	item["following_id"] = &types.AttributeValueMemberS{Value: followingID}

	err = s.dbClient.PutItem(ctx, s.dbClient.GetFollowersTableName(), item)
	if err != nil {
		return err
	}

	go func() {
		if err := s.updateFollowerTimeline(context.Background(), userID, followingID); err != nil {
			logger.LogError("Error updating follower timeline", "error", err, "follower_id", userID, "following_id", followingID)
		}
	}()

	logger.LogInfo("Follow finished successfully", "follower_id", userID, "following_id", followingID)
	return nil
}

func (s *FollowService) updateFollowerTimeline(ctx context.Context, followerID, followingID string) error {
	messages, err := s.messageService.GetUserMessages(ctx, followingID, 100)
	if err != nil {
		logger.LogError("Error getting user messages", "error", err, "following_id", followingID)
		return err
	}

	for _, message := range messages {
		if err := s.timelineService.UpdateFollowersTimeline(ctx, message); err != nil {
			continue
		}
	}

	return nil
}

func (s *FollowService) getFollowers(ctx context.Context, userID string) ([]string, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(s.dbClient.GetFollowersTableName()),
		IndexName:              aws.String("FollowingIndex"),
		KeyConditionExpression: aws.String("following_id = :following_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":following_id": &types.AttributeValueMemberS{Value: userID},
		},
	}

	result, err := s.dbClient.Query(ctx, input)
	if err != nil {
		logger.LogError("Error getting followers", "error", err, "following_id", userID)
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
