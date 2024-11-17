package services

import (
	"context"
	"time"

	"github.com/guregu/dynamo/v2"
	"github.com/the-redx/link-shortener/internal/domain"
	"github.com/the-redx/link-shortener/pkg/errs"
)

type LinkService struct {
	dynamoDB   *dynamo.DB
	linksTable dynamo.Table
}

func (s *LinkService) GetAllLinks(userId string) (*[]domain.Link, *errs.AppError) {
	var links []domain.Link

	if err := s.linksTable.Scan().Filter("'UserId' = ?", userId).All(context.TODO(), &links); err != nil {
		return nil, errs.NewUnexpectedError("Error while fetching links")
	}

	return &links, nil
}

func (s *LinkService) GetLinkByID(userId, id string) (*domain.Link, *errs.AppError) {
	var link domain.Link

	if err := s.linksTable.Get("ID", id).One(context.TODO(), &link); err != nil {
		if err == dynamo.ErrNotFound {
			return nil, errs.NewNotFoundError("Link not found")
		}

		return nil, errs.NewUnexpectedError("Error while fetching link")
	}

	if userId != link.UserId && link.Status != domain.Active {
		return nil, errs.NewNotFoundError("Link not found")
	}

	return &link, nil
}

func (s *LinkService) CreateLink(userId string, linkDTO *domain.CreateLinkDTO) (*domain.Link, *errs.AppError) {
	if linkDTO.ID == "" {
		linkDTO.ID = generateRandomID()
	}

	var link domain.Link

	if err := s.linksTable.Get("ID", linkDTO.ID).One(context.TODO(), &link); err == nil {
		return nil, errs.NewUnexpectedError("Item is already exists")
	}

	link = domain.Link{
		ID:          linkDTO.ID,
		Name:        linkDTO.Name,
		UserId:      userId,
		ShortUrl:    createShortUrlFromID(linkDTO.ID),
		Url:         linkDTO.Url,
		Status:      domain.Active,
		DateCreated: time.Now(),
		DateUpdated: time.Now(),
	}

	if err := s.linksTable.Put(link).Run(context.TODO()); err != nil {
		return nil, errs.NewUnexpectedError("Error while creating link")
	}

	return &link, nil
}

func (s *LinkService) UpdateLinkByID(userId, id string, linkDTO *domain.UpdateLinkDTO) (*domain.Link, *errs.AppError) {
	var link domain.Link

	if err := s.linksTable.Get("ID", id).Range("UserId", dynamo.Equal, userId).One(context.TODO(), &link); err != nil {
		if err == dynamo.ErrNotFound {
			return nil, errs.NewNotFoundError("Link not found")
		}

		return nil, errs.NewUnexpectedError("Error while fetching link")
	}

	name := linkDTO.Name
	if name == "" {
		name = link.Name
	}

	status := linkDTO.Status
	if status == "" {
		status = link.Status
	}

	if err := s.linksTable.Update("ID", id).Set("name", name).Set("Status", status).Set("DateUpdated", time.Now()).Run(context.TODO()); err != nil {
		return nil, errs.NewUnexpectedError("Error while updating link")
	}

	return &link, nil
}

func (s *LinkService) DeleteLinkByID(userId, id string) (*domain.Link, *errs.AppError) {
	var link domain.Link

	if err := s.linksTable.Get("ID", id).Range("UserId", dynamo.Equal, userId).One(context.TODO(), &link); err != nil {
		if err == dynamo.ErrNotFound {
			return nil, errs.NewNotFoundError("Link not found")
		}

		return nil, errs.NewUnexpectedError("Error while fetching link")
	}

	if err := s.linksTable.Delete("ID", id).Run(context.TODO()); err != nil {
		return nil, errs.NewUnexpectedError("Error while deleting link")
	}

	return &link, nil
}

func NewLinkService() LinkService {
	dynamoDB := NewDynamoDBService()
	table := dynamoDB.Table("Links")

	return LinkService{dynamoDB: dynamoDB, linksTable: table}
}
