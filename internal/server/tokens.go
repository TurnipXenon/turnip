// storage is an abstraction to s3 buckets

package server

type Token struct {
	AccessToken string
	GeneratedAt string // RFC3339
	ExpiresAt   string // RFC3339
}

type Tokens interface {
	GetOrCreateToken(ud *User) (*Token, error)
}
