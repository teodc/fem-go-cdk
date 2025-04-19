package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	UserTableName   = "FemGoCdkUsers"
	UserUsernameKey = "username"
	UserPasswordKey = "password"
)

type DynamoDBClient struct {
	session *session.Session
	client  *dynamodb.DynamoDB
}

func NewDynamoDBClient() *DynamoDBClient {
	return &DynamoDBClient{
		session: session.Must(session.NewSession()),
		client:  dynamodb.New(session.Must(session.NewSession())),
	}
}

func (db *DynamoDBClient) DoesUserExist(username string) (bool, error) {
	itemInput := &dynamodb.GetItemInput{
		TableName: aws.String(UserTableName),
		Key: map[string]*dynamodb.AttributeValue{
			UserUsernameKey: {
				S: aws.String(username),
			},
		},
	}

	res, err := db.client.GetItem(itemInput)
	if err != nil {
		return false, err
	}

	// User does not exist
	if res.Item == nil {
		return false, nil
	}

	// User exists
	return true, nil
}

func (db *DynamoDBClient) CreateUser(username, password string) error {
	// In a real app, actually hash the password
	passwordHash := password

	itemInput := &dynamodb.PutItemInput{
		TableName: aws.String(UserTableName),
		Item: map[string]*dynamodb.AttributeValue{
			UserUsernameKey: {
				S: aws.String(username),
			},
			UserPasswordKey: {
				S: aws.String(passwordHash),
			},
		},
	}

	_, err := db.client.PutItem(itemInput)
	if err != nil {
		return err
	}

	return nil
}
