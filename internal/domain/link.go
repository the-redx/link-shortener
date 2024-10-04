package domain

import (
	"time"
)

type LinkStatus string

const (
	Active          LinkStatus = "active"
	PendingApproval LinkStatus = "pending_approval"
	Suspended       LinkStatus = "suspended"
)

type Link struct {
	ID          string     `json:"id" dynamo:"ID,hash"`
	Name        string     `json:"name"`
	ShortUrl    string     `json:"shortUrl" dynamo:"-"`
	Url         string     `json:"url"`
	Status      LinkStatus `json:"status"`
	DateCreated time.Time  `json:"dateCreated"`
	DateUpdated time.Time  `json:"dateUpdated"`
	DateDeleted *time.Time `json:"dateDeleted,omitempty" dynamo:",omitempty"`
}

type CreateLinkDTO struct {
	ID   string `json:"id" validate:"max=50"`
	Name string `json:"name" validate:"max=100"`
	Url  string `json:"url" validate:"required,url,max=5000"`
}

type UpdateLinkDTO struct {
	ID   string `json:"id" validate:"max=50"`
	Name string `json:"name" validate:"min=3,max=100"`
}
