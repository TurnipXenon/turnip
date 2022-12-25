// storage is an abstraction to s3 buckets

package server

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/TurnipXenon/Turnip/internal/util"
)

const ddbTimeout = time.Second * 3

type usersDynamoDBImpl struct {
	ddb          *dynamodb.DynamoDB
	ddbTableName *string
}

var (
	UserAlreadyExists = errors.New("user already exists")
)

func NewUsersDynamoDB(d *dynamodb.DynamoDB) Users {
	s := usersDynamoDBImpl{
		ddb:          d,
		ddbTableName: aws.String("Users"),
	}
	return &s
}

// GetUser may return nil
func (u *usersDynamoDBImpl) GetUser(ud *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), ddbTimeout)
	defer cancel()

	item, err := u.ddb.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Username": {S: aws.String(ud.Username)},
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
	err = dynamodbattribute.UnmarshalMap(item.Item, &newUserData)
	if err != nil {
		return nil, util.WrapErrorWithDetails(err)
	}

	return &newUserData, nil
}

func (u *usersDynamoDBImpl) CreateUser(ctx context.Context, ud *User) error {
	// check if user already exists
	item, err := u.GetUser(ud)
	if err != nil {
		util.LogDetailedError(err)
		return util.WrapErrorWithDetails(err)
	}
	if item != nil {
		return util.WrapErrorWithDetails(UserAlreadyExists)
	}

	_, err = u.ddb.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Username":       {S: aws.String(ud.Username)},
			"HashedPassword": {S: aws.String(ud.HashedPassword)},
		},
		TableName: u.ddbTableName,
	})
	if err != nil {
		util.LogDetailedError(err)
		return util.WrapErrorWithDetails(err)
	}

	return nil
}
