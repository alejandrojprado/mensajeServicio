package database

import (
	"context"

	"mensajesService/components/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DDBClientInterface interface {
	PutItem(ctx context.Context, tableName string, item map[string]types.AttributeValue) error
	GetItem(ctx context.Context, tableName string, key map[string]types.AttributeValue) (*dynamodb.GetItemOutput, error)
	Query(ctx context.Context, input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	DeleteItem(ctx context.Context, tableName string, key map[string]types.AttributeValue) error
	GetMessagesTableName() string
	GetFollowersTableName() string
	GetTimelineTableName() string
}

type DDBClient struct {
	client              *dynamodb.Client
	tableMensajesName   string
	tableSeguidoresName string
	tableTimelineName   string
}

func NewDDBClient(ctx context.Context, cfg *config.AppConfig) (*DDBClient, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(cfg.Region))
	if err != nil {
		return nil, err
	}

	return &DDBClient{
		client:              dynamodb.NewFromConfig(awsCfg),
		tableMensajesName:   cfg.TableMensajesName,
		tableSeguidoresName: cfg.TableSeguidoresName,
		tableTimelineName:   cfg.TableTimelineName,
	}, nil
}

func (d *DDBClient) PutItem(ctx context.Context, tableName string, item map[string]types.AttributeValue) error {
	_, err := d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	return err
}

func (d *DDBClient) GetItem(ctx context.Context, tableName string, key map[string]types.AttributeValue) (*dynamodb.GetItemOutput, error) {
	return d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	})
}

func (d *DDBClient) Query(ctx context.Context, input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return d.client.Query(ctx, input)
}

func (d *DDBClient) DeleteItem(ctx context.Context, tableName string, key map[string]types.AttributeValue) error {
	_, err := d.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	})
	return err
}

func (d *DDBClient) GetMessagesTableName() string {
	return d.tableMensajesName
}

func (d *DDBClient) GetFollowersTableName() string {
	return d.tableSeguidoresName
}

func (d *DDBClient) GetTimelineTableName() string {
	return d.tableTimelineName
}
