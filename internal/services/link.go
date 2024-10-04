package services

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/guregu/dynamo/v2"
	"github.com/the-redx/link-shortener/internal/domain"
	"github.com/the-redx/link-shortener/pkg/errs"
)

type LinkService struct {
	dynamoDB   *dynamo.DB
	linksTable dynamo.Table
}

// 1. Find what is ctx
// 2. Add support for user authentication

func (s *LinkService) GetAllLinks() (*[]domain.Link, *errs.AppError) {
	var links []domain.Link

	err := s.linksTable.Scan().All(ctx, &links)
	if err != nil {
		appErr := errs.NewUnexpectedError("Error while fetching links")
		return nil, appErr
	}

	for i := 0; i < len(links); i++ {
		links[i].ShortUrl = createShortUrlFromID(links[i].ID)
	}

	return &links, nil
}

func (s *LinkService) GetLinkByID(id string) (*domain.Link, *errs.AppError) {
	var link domain.Link

	err := s.linksTable.Get("ID", id).One(ctx, &link)
	if err != nil {
		appErr := errs.NewNotFoundError("Link not found")
		return nil, appErr
	}

	link.ShortUrl = createShortUrlFromID(link.ID)

	return &link, nil
}

func (s *LinkService) CreateLink(linkDTO *domain.CreateLinkDTO) (*domain.Link, *errs.AppError) {
	var link domain.Link

	if linkDTO.ID == "" {
		linkDTO.ID = generateRandomID()
	}

	err := s.linksTable.Get("ID", linkDTO.ID).One(ctx, &link)
	if err == nil && &link != nil {
		return nil, errs.NewBadRequestError("This short URL is already in use")
	}

	link = domain.Link{
		ID:          linkDTO.ID,
		Name:        linkDTO.Name,
		ShortUrl:    createShortUrlFromID(linkDTO.ID),
		Url:         linkDTO.Url,
		Status:      "active",
		DateCreated: time.Now(),
		DateUpdated: time.Now(),
		DateDeleted: nil,
	}

	if err := s.linksTable.Put(link).Run(ctx); err != nil {
		return nil, errs.NewUnexpectedError("Error while creating link")
	}

	return &link, nil
}

func (s *LinkService) UpdateLinkByID(id string, linkDTO *domain.UpdateLinkDTO) (*domain.Link, *errs.AppError) {
	var link domain.Link

	err := s.linksTable.Get("ID", id).One(ctx, &link)
	if err != nil {
		appErr := errs.NewNotFoundError("Link not found")
		return nil, appErr
	}

	if linkDTO.ID != "" {
		link.ID = linkDTO.ID
	}

	if linkDTO.Name != "" {
		link.Name = linkDTO.Name
	}

	if err = s.linksTable.Update(id, link).Run(ctx); err != nil {
		appErr := errs.NewUnexpectedError(err.Error())
		return nil, appErr
	}

	link.ShortUrl = createShortUrlFromID(link.ID)
	return &link, nil
}

func (s *LinkService) DeleteLinkByID(id string) (*domain.Link, *errs.AppError) {
	var link domain.Link

	err := s.linksTable.Get("ID", id).One(ctx, &link)
	if err != nil {
		appErr := errs.NewNotFoundError("Link not found")
		return nil, appErr
	}

	err = s.linksTable.Delete("ID", id).Run(ctx)
	if err != nil {
		appErr := errs.NewUnexpectedError(err.Error())
		return nil, appErr
	}

	link.ShortUrl = createShortUrlFromID(link.ID)

	return &link, nil
}

func NewLinkService(awsConfig aws.Config) LinkService {
	dynamoDB := NewDynamoDBService(awsConfig)
	table := dynamoDB.Table("Links")

	return LinkService{dynamoDB: dynamoDB, linksTable: table}
}
