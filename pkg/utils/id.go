package utils

import (
	"github.com/google/uuid"
	"github.com/the-redx/link-shortener/pkg/errs"
)

func ConvertToUUID(id string) (uuid.UUID, *errs.AppError) {
	linkID, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, errs.NewBadRequestError("Invalid UUID")
	}

	return linkID, nil
}
