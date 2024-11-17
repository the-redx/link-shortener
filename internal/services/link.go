package services

import (
	"context"
	"time"

	"github.com/guregu/dynamo/v2"
	"github.com/the-redx/link-shortener/internal/domain"
	"github.com/the-redx/link-shortener/pkg/errs"
	"go.uber.org/zap"
)

type LinkService struct {
	dynamoDB   *dynamo.DB
	linksTable dynamo.Table
}

func (s *LinkService) GetAllLinks(ctx context.Context) (*[]domain.Link, *errs.AppError) {
	var links []domain.Link

	userId, ok := ctx.Value("UserID").(string)
	if !ok {
		return &links, nil
	}

	logger := ctx.Value("Logger").(*zap.Logger)

	if err := s.linksTable.Scan().Filter("'UserId' = ?", userId).All(context.TODO(), &links); err != nil {
		logger.Debug("Error while fetching links", zap.Error(err))
		return nil, errs.NewUnexpectedError("Error while fetching links")
	}

	logger.Debug("Response", zap.Any("links", links))
	return &links, nil
}

func (s *LinkService) GetLinkByID(id string, ctx context.Context) (*domain.Link, *errs.AppError) {
	var link domain.Link

	userId, ok := ctx.Value("UserID").(string)
	logger := ctx.Value("Logger").(*zap.Logger)

	if err := s.linksTable.Get("ID", id).One(context.TODO(), &link); err != nil {
		logger.Debug("Error while fetching link", zap.Error(err))

		if err == dynamo.ErrNotFound {
			return nil, errs.NewNotFoundError("Link not found")
		}

		return nil, errs.NewUnexpectedError("Error while fetching link")
	}

	logger.Debug("Result", zap.Any("link", link))

	if (!ok || userId != link.UserId) && link.Status != domain.Active {
		logger.Debug("User is not a owner and link is not active")
		return nil, errs.NewNotFoundError("Link not found")
	}

	return &link, nil
}

func (s *LinkService) CreateLink(linkDTO *domain.CreateLinkDTO, ctx context.Context) (*domain.Link, *errs.AppError) {
	var link domain.Link

	userId, ok := ctx.Value("UserID").(string)
	logger := ctx.Value("Logger").(*zap.Logger)

	if !ok {
		logger.Debug("Error while fetching user id")
		return nil, errs.NewUnexpectedError("Error while fetching user id")
	}

	if linkDTO.ID == "" {
		linkDTO.ID = generateRandomID()
	}

	if err := s.linksTable.Get("ID", linkDTO.ID).One(context.TODO(), &link); err == nil {
		logger.Debug("Item is already exists")
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

	logger.Debug("Link to create", zap.Any("link", link))

	if err := s.linksTable.Put(link).Run(context.TODO()); err != nil {
		logger.Debug("Error while creating the link", zap.Error(err))
		return nil, errs.NewUnexpectedError("Error while creating link")
	}

	logger.Debug("Link created", zap.Any("link", link))
	return &link, nil
}

func (s *LinkService) UpdateLinkByID(id string, linkDTO *domain.UpdateLinkDTO, ctx context.Context) (*domain.Link, *errs.AppError) {
	var link domain.Link

	userId, ok := ctx.Value("UserID").(string)
	logger := ctx.Value("Logger").(*zap.Logger)

	if !ok {
		logger.Debug("Error while fetching user id")
		return nil, errs.NewUnexpectedError("Error while fetching user id")
	}

	if err := s.linksTable.Get("ID", id).Range("UserId", dynamo.Equal, userId).One(context.TODO(), &link); err != nil {
		logger.Debug("Error while fetching link", zap.Error(err))
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

	logger.Debug("Link to update", zap.Any("link", link))

	if err := s.linksTable.Update("ID", id).Set("name", name).Set("Status", status).Set("DateUpdated", time.Now()).Run(context.TODO()); err != nil {
		logger.Debug("Error while updating the link", zap.Error(err))
		return nil, errs.NewUnexpectedError("Error while updating link")
	}

	logger.Debug("Link updated", zap.Any("link", link))
	return &link, nil
}

func (s *LinkService) DeleteLinkByID(id string, ctx context.Context) (*domain.Link, *errs.AppError) {
	var link domain.Link

	userId, ok := ctx.Value("UserID").(string)
	logger := ctx.Value("Logger").(*zap.Logger)

	if !ok {
		logger.Debug("Error while fetching user id")
		return nil, errs.NewUnexpectedError("Error while fetching user id")
	}

	if err := s.linksTable.Get("ID", id).Range("UserId", dynamo.Equal, userId).One(context.TODO(), &link); err != nil {
		logger.Debug("Error while fetching link", zap.Error(err))

		if err == dynamo.ErrNotFound {
			return nil, errs.NewNotFoundError("Link not found")
		}

		return nil, errs.NewUnexpectedError("Error while fetching link")
	}

	if err := s.linksTable.Delete("ID", id).Run(context.TODO()); err != nil {
		logger.Debug("Error while deleting link", zap.Error(err))
		return nil, errs.NewUnexpectedError("Error while deleting link")
	}

	logger.Debug("Link deleted", zap.Any("link", link))
	return &link, nil
}

func NewLinkService() LinkService {
	dynamoDB := NewDynamoDBService()
	table := dynamoDB.Table("Links")

	return LinkService{dynamoDB: dynamoDB, linksTable: table}
}
