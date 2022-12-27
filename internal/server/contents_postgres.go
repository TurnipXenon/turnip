package server

import (
	"context"

	"github.com/TurnipXenon/turnip/internal/clients"
	"github.com/TurnipXenon/turnip_api/rpc/turnip"
)

type contentsPostgresImpl struct {
	db        *clients.PostgresDb
	tableName string
	// todo: global secondary index
}

func (c contentsPostgresImpl) CreateContent(ctx context.Context, request *turnip.CreateContentRequest) (*turnip.Content, error) {
	//TODO implement me
	panic("implement me")
}

func (c contentsPostgresImpl) GetContentById(ctx context.Context, primary string) (*turnip.Content, error) {
	//TODO implement me
	panic("implement me")
}

func (c contentsPostgresImpl) GetAllContent(ctx context.Context) ([]*turnip.Content, error) {
	//TODO implement me
	panic("implement me")
}

func (c contentsPostgresImpl) GetContentByTag(ctx context.Context, tag string) ([]*turnip.Content, error) {
	//TODO implement me
	panic("implement me")
}

func (c contentsPostgresImpl) UpdateContent(ctx context.Context, new *turnip.Content) (*turnip.Content, error) {
	//TODO implement me
	panic("implement me")
}

func (c contentsPostgresImpl) DeleteContentById(ctx context.Context, primary string) (*turnip.Content, error) {
	//TODO implement me
	panic("implement me")
}

func NewContentsPostgres(d *clients.PostgresDb) Contents {
	// primary: primary id
	// sort: created at
	t := contentsPostgresImpl{
		db:        d,
		tableName: "Contents",
	}
	return &t
}
