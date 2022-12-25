package server

import (
	"context"
	"github.com/TurnipXenon/Turnip/internal/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/TurnipXenon/turnip_twirp/rpc/turnip"
)

type contentsDynamoDBImpl struct {
	ddb          *dynamodb.DynamoDB
	ddbTableName *string
	// todo: global secondary index
}

func NewContentsDynamoDB(d *dynamodb.DynamoDB) Contents {
	// primary: primary id
	// sort: created at
	t := contentsDynamoDBImpl{
		ddb:          d,
		ddbTableName: aws.String("Contents"),
	}
	return &t
}

func (c contentsDynamoDBImpl) CreateContent(ctx context.Context, request *turnip.CreateContentRequest) (*turnip.Content, error) {
	//TODO implement me

	// create uuid
	// vary unlikely to collide, right?
	content := turnip.Content{
		Title:         request.Title,
		Description:   request.Description,
		Content:       request.Content,
		Media:         request.Media,
		TagList:       request.TagList,
		AccessDetails: request.AccessDetails,
		Meta:          request.Meta,
		PrimaryId:     uuid.New().String(),
		CreatedAt:     timestamppb.Now(),
	}
	itemInput, err := dynamodbattribute.MarshalMap(content)
	if err != nil {
		util.LogDetailedError(err)
		return nil, err
	}

	_, err = c.ddb.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_not_exists(PrimaryId)"),
		Item:                itemInput,
		TableName:           c.ddbTableName,
	})
	if err != nil {
		util.LogDetailedError(err)
		return nil, err
	}

	// todo(turnip): tags!

	return &content, nil
}

func (c contentsDynamoDBImpl) GetContentById(ctx context.Context, primary string) (*turnip.Content, error) {
	//TODO implement me
	panic("implement me")
}

func (c contentsDynamoDBImpl) GetAllContent(ctx context.Context) ([]*turnip.Content, error) {
	//TODO implement me
	panic("implement me")
}

func (c contentsDynamoDBImpl) GetContentByTag(ctx context.Context, tag string) ([]*turnip.Content, error) {
	//TODO implement me
	panic("implement me")
}

func (c contentsDynamoDBImpl) UpdateContent(ctx context.Context, new *turnip.Content) (*turnip.Content, error) {
	//TODO implement me
	panic("implement me")
}

func (c contentsDynamoDBImpl) DeleteContentById(ctx context.Context, primary string) (*turnip.Content, error) {
	//TODO implement me
	panic("implement me")
}
