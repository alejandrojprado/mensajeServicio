package service

import (
	"context"
	"testing"
	"time"

	"mensajesService/components/database"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDDBClient struct {
	mock.Mock
}

var _ database.DDBClientInterface = (*MockDDBClient)(nil)

func (m *MockDDBClient) PutItem(ctx context.Context, tableName string, item map[string]types.AttributeValue) error {
	args := m.Called(ctx, tableName, item)
	return args.Error(0)
}

func (m *MockDDBClient) GetItem(ctx context.Context, tableName string, key map[string]types.AttributeValue) (*dynamodb.GetItemOutput, error) {
	args := m.Called(ctx, tableName, key)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (m *MockDDBClient) Query(ctx context.Context, input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*dynamodb.QueryOutput), args.Error(1)
}

func (m *MockDDBClient) GetMessagesTableName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockDDBClient) GetFollowersTableName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockDDBClient) GetTimelineTableName() string {
	args := m.Called()
	return args.String(0)
}

func TestNewMessageService(t *testing.T) {
	mockDB := &MockDDBClient{}
	service := NewMessageService(mockDB)

	assert.NotNil(t, service)
	assert.Equal(t, mockDB, service.dbClient)
}

func TestCreateMessage_Success(t *testing.T) {
	mockDB := &MockDDBClient{}
	service := NewMessageService(mockDB)

	ctx := context.Background()
	userID := "user123"
	content := "Test message content"

	mockDB.On("GetMessagesTableName").Return("messages-table")
	mockDB.On("PutItem", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("map[string]types.AttributeValue")).Return(nil)
	mockDB.On("GetFollowersTableName").Return("followers-table")
	mockDB.On("Query", mock.Anything, mock.AnythingOfType("*dynamodb.QueryInput")).Return(&dynamodb.QueryOutput{
		Items: []map[string]types.AttributeValue{},
	}, nil)

	message, err := service.CreateMessage(ctx, userID, content)

	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, userID, message.UserID)
	assert.Equal(t, content, message.Content)
	assert.NotEmpty(t, message.ID)
	assert.True(t, time.Since(message.CreatedAt) < time.Second)

	// wait time (It could be better to avoid the fixed wait time)
	time.Sleep(100 * time.Millisecond)

	mockDB.AssertNumberOfCalls(t, "GetMessagesTableName", 1)
	mockDB.AssertNumberOfCalls(t, "PutItem", 1)
}

func TestCreateMessage_DatabaseError(t *testing.T) {
	mockDB := &MockDDBClient{}
	service := NewMessageService(mockDB)

	ctx := context.Background()
	userID := "user123"
	content := "Test message content"

	mockDB.On("GetMessagesTableName").Return("messages-table")
	mockDB.On("PutItem", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("map[string]types.AttributeValue")).Return(assert.AnError)

	message, err := service.CreateMessage(ctx, userID, content)

	assert.Error(t, err)
	assert.Nil(t, message)

	mockDB.AssertExpectations(t)
}

func TestGetUserMessages_Success(t *testing.T) {
	mockDB := &MockDDBClient{}
	service := NewMessageService(mockDB)

	ctx := context.Background()
	userID := "user123"
	limit := 10

	messages := []map[string]types.AttributeValue{
		{
			"message_id": &types.AttributeValueMemberS{Value: "msg1"},
			"user_id":    &types.AttributeValueMemberS{Value: userID},
			"content":    &types.AttributeValueMemberS{Value: "Test content 1"},
			"created_at": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
		},
	}

	mockDB.On("GetMessagesTableName").Return("messages-table")
	mockDB.On("Query", ctx, mock.AnythingOfType("*dynamodb.QueryInput")).Return(&dynamodb.QueryOutput{
		Items: messages,
	}, nil)

	result, err := service.GetUserMessages(ctx, userID, limit)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "msg1", result[0].ID)
	assert.Equal(t, "Test content 1", result[0].Content)

	mockDB.AssertExpectations(t)
}
