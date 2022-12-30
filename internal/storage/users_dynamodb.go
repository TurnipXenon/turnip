// storage is an abstraction to s3 buckets

package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/TurnipXenon/turnip/internal/util"
)

const ddbTimeout = time.Second * 3

type usersDynamoDBImpl struct {
	ddb          *dynamodb.Client
	ddbTableName *string
}

func NewUsersDynamoDB(d *dynamodb.Client) Users {
	s := usersDynamoDBImpl{
		ddb:          d,
		ddbTableName: aws.String("Users"),
	}
	return &s
}

// GetUser may return nil
func (u *usersDynamoDBImpl) GetUser(_ context.Context, ud *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), ddbTimeout)
	defer cancel()

	item, err := u.ddb.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"Username": &types.AttributeValueMemberS{Value: ud.Username},
		},
		TableName: u.ddbTableName,
	})
	if err != nil {
		return nil, util.WrapErrorWithDetails(err)
	}

	if item.Item == nil {
		return nil, nil
	}

	newUserData := User{}
	err = attributevalue.UnmarshalMap(item.Item, &newUserData)
	if err != nil {
		return nil, util.WrapErrorWithDetails(err)
	}

	return &newUserData, nil
}

func (u *usersDynamoDBImpl) CreateUser(ctx context.Context, ud *User) error {
	// check if migration already exists
	item, err := u.GetUser(ctx, ud)
	if err != nil {
		util.LogDetailedError(err)
		return util.WrapErrorWithDetails(err)
	}
	if item != nil {
		return util.WrapErrorWithDetails(UserAlreadyExists)
	}

	_, err = u.ddb.PutItem(ctx, &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"Username":       &types.AttributeValueMemberS{Value: ud.Username},
			"HashedPassword": &types.AttributeValueMemberS{Value: ud.HashedPassword},
		},
		TableName: u.ddbTableName,
	})
	if err != nil {
		util.LogDetailedError(err)
		return util.WrapErrorWithDetails(err)
	}

	return nil
}
