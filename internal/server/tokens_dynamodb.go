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
	// todo: global secondary index
}

func NewTokensDynamoDB(d *dynamodb.DynamoDB) Tokens {
	t := tokensDynamoDBImpl{
		ddb:          d,
		ddbTableName: aws.String("Tokens"),
	}
	return &t
}

func (t *tokensDynamoDBImpl) GetOrCreateTokenByUsername(ctx context.Context, ud *User) (*models.Token, error) {
	token := models.Token{}

	// (1) if token exists
	query, err := t.ddb.QueryWithContext(ctx, &dynamodb.QueryInput{
		KeyConditions: map[string]*dynamodb.Condition{
			"Username": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(ud.Username),
					},
				},
			},
		},
		IndexName: aws.String("UsernameIndex"),
		TableName: t.ddbTableName,
	})
	if err != nil {
		util.LogDetailedError(err)
		return nil, err
	}
	if len(query.Items) > 0 {
		// it exists!
		err = dynamodbattribute.UnmarshalMap(query.Items[0], &token)
		if err != nil {
			util.LogDetailedError(err)
			return nil, err
		}

		return &token, nil
	}

	// (2) if token does not exist
	token.AccessToken, err = generateSecureToken(40)
	if err != nil {
		util.LogDetailedError(err)
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
		util.LogDetailedError(err)
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

func (t *tokensDynamoDBImpl) GetToken(accessToken string) (*models.Token, error) {
	// todo: how do we make this inline with twirp generated server stubs?
	ctx, cancel := context.WithTimeout(context.TODO(), ddbTimeout)
	defer cancel()

	item, err := t.ddb.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"TokenAccess": {S: aws.String(accessToken)},
		},
		TableName: t.ddbTableName,
	})
	if err != nil {
		util.LogDetailedError(err)
		return nil, err
	}
	if item.Item != nil {
		token := models.Token{}
		err = dynamodbattribute.UnmarshalMap(item.Item, &token)
		if err != nil {
			util.LogDetailedError(err)
			return nil, err
		}

		return &token, nil
	}

	return nil, nil
}
