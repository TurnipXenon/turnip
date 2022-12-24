package models

type Token struct {
	AccessToken string
	GeneratedAt string // RFC3339
	ExpiresAt   string // RFC3339
}
