// storage is an abstraction to s3 buckets

package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/TurnipXenon/Turnip/internal/util"
	"github.com/TurnipXenon/Turnip/pkg/models"
)

type tokensDynamoDBImpl struct {
	ddb          *dynamodb.DynamoDB
	ddbTableName *string
}

func NewTokensDynamoDB(d *dynamodb.DynamoDB) Tokens {
	t := tokensDynamoDBImpl{
		ddb:          d,
		ddbTableName: aws.String("Tokens"),
	}
	return &t
}

func (t *tokensDynamoDBImpl) GetOrCreateToken(ud *User) (*models.Token, error) {
	// todo: add this pattern to all calls here???
	ctx, cancel := context.WithTimeout(context.TODO(), ddbTimeout)
	defer cancel()

	token := models.Token{}

	// (1) if token exists
	item, err := t.ddb.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Username": {S: aws.String(ud.Username)},
		},
		TableName: t.ddbTableName,
	})
	if err != nil {
		util.LogDetailedError(err)
		return nil, err
	}
	if item.Item != nil {
		// it exists!
		err = dynamodbattribute.UnmarshalMap(item.Item, &token)
		if err != nil {
			util.LogDetailedError(err)
			return nil, err
		}

		return &token, nil
	}

	// (2) if token does not exist
	token.AccessToken, err = generateSecureToken(40)
	if err != nil {
		fmt.Printf("GetOrCreateToken: Error: %s\n", err.Error())
		return nil, err
	}

	//dt := time.Now()
	//token.GeneratedAt = dt.Format(time.RFC3339)
	//expiryTime := time.Now().Add(time.Hour * 24)
	//token.ExpiresAt = expiryTime.Format(time.RFC3339)

	_, err = t.ddb.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Username":    {S: aws.String(ud.Username)},
			"AccessToken": {S: aws.String(token.AccessToken)}, // todo: we could support expirable tokens later...
		},
		TableName: t.ddbTableName,
	})
	if err != nil {
		fmt.Printf("GetOrCreateToken: Error: %s\n", err.Error())
		return nil, err
	}
	return &token, err
}

func generateSecureToken(length int) (string, error) {
	// from Andzej Maciusovic @ https://stackoverflow.com/a/59457748/17836168
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		fmt.Printf("generateSecureToken: Error: %s\n", err.Error())
		return "", err
	}
	return hex.EncodeToString(b), nil
}
