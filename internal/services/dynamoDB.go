package services

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/guregu/dynamo/v2"
	"github.com/the-redx/link-shortener/pkg/utils"
)

func NewDynamoDBService() *dynamo.DB {
	dynamoEndpoint := os.Getenv("DYNAMODB_ENDPOINT")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: dynamoEndpoint}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)

	if err != nil {
		log.Fatal("Error loading AWS config")
	}

	utils.Logger.Debug("DynamoDB config loaded", dynamoEndpoint)

	return dynamo.New(cfg)
}

func GetOrCreateTable(db *dynamo.DB, tableName string, from interface{}) dynamo.Table {
	if err := db.CreateTable(tableName, from).Run(context.TODO()); err != nil {
		utils.Logger.Info(err)
	}

	return db.Table(tableName)
}
