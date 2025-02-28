package service

import (
	"crypto/sha256"
	"net/url"
	"strings"
)

type HTPPModel struct {
	URL string `json:"url"`
}

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"

func ShortURL(input string) string {
	hash := sha256.Sum256([]byte(input))
	var builder strings.Builder
	builder.Grow(10)
	for i := 0; i < 10; i++ {
		index := int(hash[i]) % len(alphabet)
		builder.WriteByte(alphabet[index])
	}
	return builder.String()
}

func IsValidURL(input string) bool {
	parsed_URL, err := url.ParseRequestURI(input)
	if err != nil {
		return false
	}
	return parsed_URL.Scheme != "" && parsed_URL.Host != ""
}
