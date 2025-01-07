package services

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/guregu/dynamo/v2"
	"github.com/rs/xid"
	"github.com/the-redx/link-shortener/internal/domain"
	"github.com/the-redx/link-shortener/pkg/errs"
	"github.com/the-redx/link-shortener/pkg/utils"
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

	logger := ctx.Value("Logger").(*zap.SugaredLogger)

	if err := s.linksTable.Scan().Filter("'UserId' = ? AND 'Status' = ?", userId, "active").All(context.TODO(), &links); err != nil {
		logger.Debug("Error while fetching links", zap.Error(err))
		return nil, errs.NewUnexpectedError("Error while fetching links")
	}

	logger.Debug("Response", zap.Any("links", links))
	return &links, nil
}

func (s *LinkService) GetLinkByID(id string, ctx context.Context) (*domain.Link, *errs.AppError) {
	userId, ok := ctx.Value("UserID").(string)
	logger := ctx.Value("Logger").(*zap.SugaredLogger)

	link, appErr := s.getLinkByID(id, ctx)
	if appErr != nil {
		return nil, appErr
	}

	if !ok || userId != link.UserId {
		logger.Debug("User is not a owner and link is not active")
		return nil, errs.NewForbiddenError("You don't have access to this link")
	}

	return link, nil
}

func (s *LinkService) GetLinkByIDForRedirect(id string, ctx context.Context) (*domain.Link, *errs.AppError) {
	logger := ctx.Value("Logger").(*zap.SugaredLogger)

	link, appErr := s.getLinkByID(id, ctx)
	if appErr != nil {
		return nil, appErr
	}

	logger.Debug("Link status", zap.Any("Status", link.Status))

	if link.Status != domain.Active {
		logger.Debug("Link is not active")
		return nil, errs.NewNotFoundError("Link not found")
	}

	// Increment the Redirects counter
	if err := s.linksTable.Update("ID", id).Set("Redirects", link.Redirects+1).Run(context.TODO()); err != nil {
		logger.Debug("Error while updating the link", zap.Error(err))
	}

	return link, nil
}

func (s *LinkService) CreateLink(linkDTO *domain.CreateLinkDTO, ctx context.Context) (*domain.Link, *errs.AppError) {
	var link domain.Link

	userId, ok := ctx.Value("UserID").(string)
	logger := ctx.Value("Logger").(*zap.SugaredLogger)

	if !ok {
		logger.Debug("Error while fetching user id")
		return nil, errs.NewUnexpectedError("Error while fetching user id")
	}

	re := regexp.MustCompile(`[^\d\p{Latin}- ]`)
	linkID := re.ReplaceAllString(linkDTO.ID, "")
	linkID = strings.ReplaceAll(strings.Trim(linkID, " "), " ", "-")

	if linkID == "" {
		linkID = xid.New().String()
		utils.Logger.Debug("Empty link ID. Use generated ID")
	}

	utils.Logger.Debug(zap.String("linkID", linkID))

	if err := s.linksTable.Get("ID", linkID).One(context.TODO(), &link); err == nil {
		logger.Debug("Item is already exists")
		return nil, errs.NewUnexpectedError("Item is already exists")
	}

	link = domain.Link{
		ID:          linkID,
		Name:        linkDTO.Name,
		UserId:      userId,
		ShortUrl:    createShortUrlFromID(linkID),
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
	userId := ctx.Value("UserID").(string)
	logger := ctx.Value("Logger").(*zap.SugaredLogger)

	link, appErr := s.getLinkByID(id, ctx)
	if appErr != nil {
		return nil, appErr
	}

	if link.UserId != userId {
		logger.Debug("User is not a owner")
		return nil, errs.NewForbiddenError("You don't have access to this link")
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

	if err := s.linksTable.Update("ID", id).Set("Name", name).Set("Status", status).Set("DateUpdated", time.Now().UTC().Format(time.RFC3339)).Run(context.TODO()); err != nil {
		logger.Debug("Error while updating the link", zap.Error(err))
		return nil, errs.NewUnexpectedError("Error while updating link")
	}

	logger.Debug("Link updated", zap.Any("link", link))

	link, appErr = s.getLinkByID(id, ctx)
	if appErr != nil {
		return nil, appErr
	}

	return link, nil
}

func (s *LinkService) DeleteLinkByID(id string, ctx context.Context) (*domain.Link, *errs.AppError) {
	userId := ctx.Value("UserID").(string)
	logger := ctx.Value("Logger").(*zap.SugaredLogger)

	link, appErr := s.getLinkByID(id, ctx)
	if appErr != nil {
		return nil, appErr
	}

	if link.UserId != userId {
		logger.Debug("User is not a owner")
		return nil, errs.NewForbiddenError("You don't have access to this link")
	}

	if err := s.linksTable.Delete("ID", id).Run(context.TODO()); err != nil {
		logger.Debug("Error while deleting link", zap.Error(err))
		return nil, errs.NewUnexpectedError("Error while deleting link")
	}

	logger.Debug("Link deleted", zap.Any("link", link))
	return link, nil
}

func (s *LinkService) getLinkByID(id string, ctx context.Context) (*domain.Link, *errs.AppError) {
	logger := ctx.Value("Logger").(*zap.SugaredLogger)
	var link domain.Link

	if err := s.linksTable.Get("ID", id).One(context.TODO(), &link); err != nil {
		if err == dynamo.ErrNotFound {
			logger.Debug("Link not found")
			return nil, errs.NewNotFoundError("Link not found")
		}

		logger.Debug("Error while fetching link", zap.Error(err))
		return nil, errs.NewUnexpectedError("Error while fetching link")
	}

	logger.Debug("Link fetched", zap.Any("link", link))
	return &link, nil
}

func NewLinkService() LinkService {
	dynamoDB := NewDynamoDBService()
	table := GetOrCreateTable(dynamoDB, "Links", domain.Link{})

	return LinkService{dynamoDB: dynamoDB, linksTable: table}
}
