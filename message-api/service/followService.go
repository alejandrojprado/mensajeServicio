package service

import (
	"context"
	"time"

	"mensajesService/components/database"
	"mensajesService/message-api/model"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type FollowServiceInterface interface {
	FollowUser(ctx context.Context, userID, followingID string) error
}

type FollowService struct {
	dbClient database.DDBClientInterface
}

func NewFollowService(dbClient database.DDBClientInterface) *FollowService {
	return &FollowService{dbClient: dbClient}
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

	item["PK"] = &types.AttributeValueMemberS{Value: "FOLLOWER#" + userID}
	item["SK"] = &types.AttributeValueMemberS{Value: "FOLLOWING#" + followingID}

	err = s.dbClient.PutItem(ctx, s.dbClient.GetFollowersTableName(), item)
	if err != nil {
		return err
	}

	return nil
}
