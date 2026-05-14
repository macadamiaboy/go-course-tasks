package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyHeader   = errors.New("authorization header is empty")
	ErrInvalidFormat = errors.New("invalid authorization header format")
	ErrEmptyToken    = errors.New("token is empty")
)

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", ErrEmptyHeader
	}
	token, ok := strings.CutPrefix(header, "Bearer ")
	if !ok {
		return "", ErrInvalidFormat
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return "", ErrEmptyToken
	}
	return token, nil
}

func main() {
	cases := []struct {
		header string
	}{
		{"Bearer my-secret-token"},
		{""},
		{"Token my-secret-token"},
		{"Bearer "},
	}

	for _, c := range cases {
		token, err := extractBearerToken(c.header)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Println("extracted:", token)
		}
	}
}
