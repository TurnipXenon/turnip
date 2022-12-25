// storage is an abstraction to s3 buckets

package server

import (
	"context"
	"github.com/TurnipXenon/turnip_twirp/rpc/turnip"
)

type Contents interface {
	CreateContent(ctx context.Context, request *turnip.CreateContentRequest) (*turnip.Content, error)
	GetContentById(ctx context.Context, primary string) (*turnip.Content, error)
	GetAllContent(ctx context.Context) ([]*turnip.Content, error)
	GetContentByTag(ctx context.Context, tag string) ([]*turnip.Content, error)

	// UpdateContent returns old version
	UpdateContent(ctx context.Context, new *turnip.Content) (*turnip.Content, error)
	DeleteContentById(ctx context.Context, primary string) (*turnip.Content, error)
}
