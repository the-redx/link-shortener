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
	ID          string     `json:"id" dynamo:",hash"`
	Name        string     `json:"name" dynamo:"Name"`
	UserId      string     `json:"userId" dynamo:"UserId"`
	ShortUrl    string     `json:"shortUrl" dynamo:"ShortUrl"`
	Url         string     `json:"url" dynamo:"Url"`
	Status      LinkStatus `json:"status" dynamo:"Status"`
	DateCreated time.Time  `json:"dateCreated" dynamo:",unixtime"`
	DateUpdated time.Time  `json:"dateUpdated" dynamo:",unixtime"`
}

type CreateLinkDTO struct {
	ID   string `json:"id" validate:"max=50"`
	Name string `json:"name" validate:"min=3,max=100"`
	Url  string `json:"url" validate:"required,url,max=5000"`
}

type UpdateLinkDTO struct {
	Name   string     `json:"name" validate:"min=3,max=100"`
	Status LinkStatus `json:"status"`
}
