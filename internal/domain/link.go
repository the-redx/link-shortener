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
	ID          string     `json:"id" dynamodbav:"ID"`
	Name        string     `json:"name" dynamodbav:"Name"`
	ShortUrl    string     `json:"shortUrl" dynamodbav:"ShortUrl"`
	Url         string     `json:"url" dynamodbav:"Url"`
	Status      LinkStatus `json:"status" dynamodbav:"Status"`
	DateCreated time.Time  `json:"dateCreated" dynamodbav:"DateCreated"`
	DateUpdated time.Time  `json:"dateUpdated" dynamodbav:"DateUpdated"`
	DateDeleted *time.Time `json:"dateDeleted,omitempty" dynamodbav:"DateDeleted,omitempty"`
}

type CreateLinkDTO struct {
	ID   string `json:"id" validate:"max=50"`
	Name string `json:"name" validate:"max=100"`
	Url  string `json:"url" validate:"required,url,max=5000"`
}

type UpdateLinkDTO struct {
	Name   string     `json:"name" validate:"min=3,max=100"`
	Url    string     `json:"url" validate:"url,max=5000"`
	Status LinkStatus `json:"status"`
}
