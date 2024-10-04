package services

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/guregu/dynamo/v2"
)

func NewDynamoDBService(cfg aws.Config) *dynamo.DB {
	db := dynamo.New(cfg)
	return db
}
