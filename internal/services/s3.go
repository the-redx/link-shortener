package services

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/the-redx/link-shortener/pkg/utils"
)

const LINK_ATTACHMENTS_BUCKET = "links-attachments"

func NewS3Service() *s3.Client {
	environment := os.Getenv("APP_ENV")
	s3Endpoint := os.Getenv("S3_ENDPOINT")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-north-1"))
	if err != nil {
		utils.Logger.Fatal("Error loading AWS config")
	}

	utils.Logger.Debug("S3 config created")

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		if environment == "development" && s3Endpoint != "" {
			utils.Logger.Debugf("S3 base endpoint is set to %s", s3Endpoint)
			o.UsePathStyle = true
			o.BaseEndpoint = &s3Endpoint
		}
	})
}
