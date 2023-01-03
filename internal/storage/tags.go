package storage

import (
	"context"

	"github.com/TurnipXenon/turnip_api/rpc/turnip"
)

type Tags interface {
	UpdateTags(ctx context.Context, content *turnip.Content) error
	GetTagsByContent(ctx context.Context, content *turnip.Content) ([]string, error)
	GetContentIdsByTag(ctx context.Context, tagList []string) ([]string, error)
	DeleteTagsByContentId(ctx context.Context, primaryId string) error
}
