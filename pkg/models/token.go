package models

type Token struct {
	AccessToken string // hash key
	Username    string // sort key
	GeneratedAt string // RFC3339
	ExpiresAt   string // RFC3339
}
