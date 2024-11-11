package services

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/the-redx/link-shortener/internal/domain"
	"github.com/the-redx/link-shortener/pkg/errs"
)

type LinkService struct {
	dynamoDB *dynamodb.Client
}

// 1. Add support for user authentication
// 2. Fix update link method

func (s *LinkService) GetAllLinks() (*[]domain.Link, *errs.AppError) {
	input := &dynamodb.ScanInput{TableName: aws.String("Links")}

	result, err := s.dynamoDB.Scan(context.TODO(), input)
	if err != nil {
		appErr := errs.NewUnexpectedError("Error while fetching links")
		return nil, appErr
	}

	var links []domain.Link
	for _, item := range result.Items {
		var link domain.Link
		err = attributevalue.UnmarshalMap(item, &link)

		if err == nil {
			links = append(links, link)
		}
	}

	return &links, nil
}

func (s *LinkService) GetLinkByID(id string) (*domain.Link, *errs.AppError) {
	result, err := s.getItemById(id)
	if err != nil {
		return nil, errs.NewUnexpectedError("Error while fetching link")
	}

	if result == nil {
		return nil, errs.NewNotFoundError("Link not found")
	}

	return result, nil
}

func (s *LinkService) CreateLink(linkDTO *domain.CreateLinkDTO) (*domain.Link, *errs.AppError) {
	if linkDTO.ID == "" {
		linkDTO.ID = generateRandomID()
	}

	result, err := s.getItemById(linkDTO.ID)
	if err != nil {
		return nil, errs.NewUnexpectedError("Error while creating link")
	}

	if result != nil {
		return nil, errs.NewBadRequestError("This short URL is already in use")
	}

	link := domain.Link{
		ID:          linkDTO.ID,
		Name:        linkDTO.Name,
		ShortUrl:    createShortUrlFromID(linkDTO.ID),
		Url:         linkDTO.Url,
		Status:      "active",
		DateCreated: time.Now(),
		DateUpdated: time.Now(),
		DateDeleted: nil,
	}

	item, err := attributevalue.MarshalMap(link)
	if err != nil {
		return nil, errs.NewUnexpectedError("Error while creating link")
	}

	putInput := &dynamodb.PutItemInput{
		TableName: aws.String("Links"),
		Item:      item,
	}

	if _, err := s.dynamoDB.PutItem(context.TODO(), putInput); err != nil {
		return nil, errs.NewUnexpectedError("Error while creating link")
	}

	return &link, nil
}

func (s *LinkService) UpdateLinkByID(id string, linkDTO *domain.UpdateLinkDTO) (*domain.Link, *errs.AppError) {
	result, err := s.getItemById(id)
	if err != nil {
		return nil, errs.NewUnexpectedError("Error while updating link")
	}

	if result == nil {
		return nil, errs.NewNotFoundError("Link not found")
	}

	name := linkDTO.Name
	if name == "" {
		name = result.Name
	}

	url := linkDTO.Url
	if url == "" {
		url = result.Url
	}

	status := linkDTO.Status
	if status == "" {
		status = result.Status
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("Links"),
		Key: map[string]dynamodbTypes.AttributeValue{
			"ID": &dynamodbTypes.AttributeValueMemberS{Value: id},
		},
		UpdateExpression: aws.String("SET Name = :name, Url = :url, Status = :status, DateUpdated = :dateUpdated"),
		ExpressionAttributeValues: map[string]dynamodbTypes.AttributeValue{
			":name":        &dynamodbTypes.AttributeValueMemberN{Value: name},
			":url":         &dynamodbTypes.AttributeValueMemberN{Value: url},
			":status":      &dynamodbTypes.AttributeValueMemberS{Value: fmt.Sprintf("%s", status)},
			":dateUpdated": &dynamodbTypes.AttributeValueMemberS{Value: time.Now().String()},
		},
		ReturnValues: "ALL_NEW",
	}

	updatedRes, err := s.dynamoDB.UpdateItem(context.TODO(), input)
	if err != nil {
		return nil, errs.NewUnexpectedError("Error while updating link")
	}

	var link domain.Link
	err = attributevalue.UnmarshalMap(updatedRes.Attributes, &link)
	if err != nil {
		return nil, errs.NewUnexpectedError("Error while updating link")
	}

	return &link, nil
}

func (s *LinkService) DeleteLinkByID(id string) (*domain.Link, *errs.AppError) {
	var link domain.Link

	result, err := s.getItemById(id)
	if err != nil {
		return nil, errs.NewUnexpectedError("Error while deleting link")
	}

	if result == nil {
		return nil, errs.NewNotFoundError("Link not found")
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("Links"),
		Key: map[string]dynamodbTypes.AttributeValue{
			"ID": &dynamodbTypes.AttributeValueMemberS{Value: id},
		},
	}

	if _, err := s.dynamoDB.DeleteItem(context.TODO(), input); err != nil {
		appErr := errs.NewUnexpectedError("Error while deleting link")
		return nil, appErr
	}

	return &link, nil
}

func (s *LinkService) getItemById(id string) (*domain.Link, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Links"),
		Key: map[string]dynamodbTypes.AttributeValue{
			"ID": &dynamodbTypes.AttributeValueMemberS{Value: id},
		},
	}

	result, err := s.dynamoDB.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	var link domain.Link
	err = attributevalue.UnmarshalMap(result.Item, &link)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func NewLinkService() LinkService {
	dynamoDB := NewDynamoDBService()

	return LinkService{dynamoDB: dynamoDB}
}
