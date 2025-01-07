package services

import (
	"os"
)

func createShortUrlFromID(id string) string {
	domainName := os.Getenv("DOMAIN_NAME")

	if domainName == "" {
		return id
	}

	return domainName + "/" + id
}
