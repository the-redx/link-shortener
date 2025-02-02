package services

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/guregu/dynamo/v2"
	"github.com/the-redx/link-shortener/pkg/utils"
)

func NewDynamoDBService() *dynamo.DB {
	environment := os.Getenv("APP_ENV")
	dynamoEndpoint := os.Getenv("DYNAMODB_ENDPOINT")

	utils.Logger.Debugf("Creating DynamoDB config for %s environment. AWS_ACCESS_KEY_ID = %s, AWS_SECRET_ACCESS_KEY = %s", environment, os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"))

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-north-1"))
	if err != nil {
		utils.Logger.Fatal("Error loading AWS config")
	}

	utils.Logger.Debug("DynamoDB config created")

	return dynamo.New(cfg, func(o *dynamodb.Options) {
		if environment == "development" && dynamoEndpoint != "" {
			utils.Logger.Debugf("DynamoDB base endpoint is set to %s", dynamoEndpoint)
			o.BaseEndpoint = &dynamoEndpoint
		}
	})
}

func GetOrCreateTable(db *dynamo.DB, tableName string, from interface{}) dynamo.Table {
	utils.Logger.Debugf("Checking if table %s exists", tableName)

	tables, err := db.ListTables().All(context.TODO())
	if err != nil {
		utils.Logger.Panic(err)
	}

	utils.Logger.Debugf("Tables: %v", tables)

	for _, table := range tables {
		if table == tableName {
			utils.Logger.Debug("Table found")
			return db.Table(tableName)
		}
	}

	utils.Logger.Debug("Table not found. Creating...")

	if err := db.CreateTable(tableName, from).Run(context.TODO()); err != nil {
		utils.Logger.Panic(err)
	}

	utils.Logger.Debug("Table created")

	return db.Table(tableName)
}
