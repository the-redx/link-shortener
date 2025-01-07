package domain

import (
	"time"
)

type LinkStatus string

const (
	Active LinkStatus = "active"
	Paused LinkStatus = "paused"
)

type Link struct {
	ID          string     `json:"id" dynamo:"ID,hash"`
	Name        string     `json:"name" dynamo:"Name"`
	UserId      string     `json:"-" dynamo:"UserId"`
	ShortUrl    string     `json:"shortUrl" dynamo:"ShortUrl"`
	Redirects   int        `json:"redirects" dynamo:"Redirects"`
	Url         string     `json:"url" dynamo:"Url"`
	Status      LinkStatus `json:"status" dynamo:"Status"`
	DateCreated time.Time  `json:"dateCreated" dynamo:"DateCreated,unixtime"`
	DateUpdated time.Time  `json:"dateUpdated" dynamo:"DateUpdated,unixtime"`
}

type CreateLinkDTO struct {
	ID   string `json:"id" validate:"max=30"`
	Name string `json:"name" validate:"min=3,max=100"`
	Url  string `json:"url" validate:"required,url,max=5000"`
}

type UpdateLinkDTO struct {
	Name   string     `json:"name" validate:"min=3,max=100"`
	Status LinkStatus `json:"status" validate:"oneof=active paused"`
}
