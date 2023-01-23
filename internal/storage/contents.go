// storage is an abstraction to s3 buckets

package storage

import (
	"context"

	"github.com/TurnipXenon/turnip_api/rpc/turnip"
)

type Contents interface {
	CreateContent(ctx context.Context, request *turnip.ContentRequestResponse, user *turnip.User) (*turnip.Content, error)
	GetContentById(ctx context.Context, primary string) (*turnip.Content, error)
	GetContentBySlug(ctx context.Context, slug string) (*turnip.Content, error)
	GetAllContent(ctx context.Context) ([]*turnip.Content, error)
	GetContentByTagInclusive(ctx context.Context, tag []string) ([]*turnip.Content, error)
	GetContentByTagStrict(ctx context.Context, tag []string) ([]*turnip.Content, error)

	// UpdateContent returns old version?
	UpdateContent(ctx context.Context, new *turnip.Content) (*turnip.Content, error)
	DeleteContentById(ctx context.Context, primary string) (*turnip.Content, error)

	// todo: batch get?
}
