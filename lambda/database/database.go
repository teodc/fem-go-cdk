package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"lambda/types"
)

const (
	UserTableName   = "FemGoCdkUsers"
	UserUsernameKey = "username"
	UserPasswordKey = "password"
)

type UserStore interface {
	DoesUserExist(username string) (bool, error)
	PersistUser(user *types.User) error
}

type DynamoDBStore struct {
	session *session.Session
	client  *dynamodb.DynamoDB
}

func NewDynamoDBStore() *DynamoDBStore {
	return &DynamoDBStore{
		session: session.Must(session.NewSession()),
		client:  dynamodb.New(session.Must(session.NewSession())),
	}
}

func (store *DynamoDBStore) DoesUserExist(username string) (bool, error) {
	itemInput := &dynamodb.GetItemInput{
		TableName: aws.String(UserTableName),
		Key: map[string]*dynamodb.AttributeValue{
			UserUsernameKey: {
				S: aws.String(username),
			},
		},
	}

	res, err := store.client.GetItem(itemInput)
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

func (store *DynamoDBStore) PersistUser(user *types.User) error {
	itemInput := &dynamodb.PutItemInput{
		TableName: aws.String(UserTableName),
		Item: map[string]*dynamodb.AttributeValue{
			UserUsernameKey: {
				S: aws.String(user.Username),
			},
			UserPasswordKey: {
				S: aws.String(user.PasswordHash),
			},
		},
	}

	_, err := store.client.PutItem(itemInput)
	if err != nil {
		return err
	}

	return nil
}
