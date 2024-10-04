package services

import (
	"crypto/rand"
	"fmt"
	"os"

	"github.com/itchyny/base58-go"
)

func generateRandomID() string {
	randomBytesLength := 10
	shortLinkLength := 6

	// Generate random bytes
	randomBytes := make([]byte, randomBytesLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Println("Error generating random bytes:", err)
		return ""
	}

	encoded, err := base58.BitcoinEncoding.Encode(randomBytes)
	if err != nil {
		fmt.Println("Error encoding to base58:", err)
		return ""
	}

	shortLink := string(encoded)
	if len(shortLink) > shortLinkLength {
		return shortLink[:shortLinkLength]
	}
	return shortLink
}

func createShortUrlFromID(id string) string {
	domainName := os.Getenv("DOMAIN_NAME")

	if domainName == "" {
		return id
	}

	return domainName + "/" + id
}
