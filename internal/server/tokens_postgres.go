// storage is an abstraction to s3 buckets

package server

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/TurnipXenon/turnip_api/rpc/turnip"

	"github.com/TurnipXenon/turnip/internal/clients"
	"github.com/TurnipXenon/turnip/internal/server/sql/migration"
	"github.com/TurnipXenon/turnip/internal/util"
)

type tokensPostgresImpl struct {
	db          *clients.PostgresDb
	dbTableName string
	// todo: global secondary index
}

func (t *tokensPostgresImpl) GetTableName() string {
	return t.dbTableName
}

func (t *tokensPostgresImpl) GetMigrationSequence() []migration.Migration {
	return []migration.Migration{
		migration.NewGenericMigration(migration.MigrateToken0001),
	}
}

func (t *tokensPostgresImpl) GetOrCreateTokenByUsername(ctx context.Context, ud *User) (*turnip.Token, error) {
	token := turnip.Token{}

	// (1) check if token exists
	// todo: optimize binary output for timestamp
	var timeCreatedAt, timeExpiresAt pgtype.Timestamp
	row := t.db.Pool.QueryRow(ctx, `SELECT *
FROM "Token" t
WHERE username=$1
LIMIT 1`, ud.Username)
	err := row.Scan(&token.AccessToken, &token.Username, &timeCreatedAt, &timeExpiresAt)
	if err != nil && err != pgx.ErrNoRows {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
	}
	if err == nil {
		if timeCreatedAt.Valid {
			token.CreatedAt = timestamppb.New(timeCreatedAt.Time)
		}
		if timeExpiresAt.Valid {
			token.ExpiresAt = timestamppb.New(timeExpiresAt.Time)
		}

		return &token, nil
	}

	// (2) if token does not exist
	token.Username = ud.Username
	token.AccessToken = uuid.New().String()

	dt := time.Now()
	token.CreatedAt = timestamppb.New(dt)
	expiryTime := time.Now().Add(time.Hour * 24)
	token.ExpiresAt = timestamppb.New(expiryTime)

	//INSERT INTO public."Token" (access_token, username, created_at, expires_at) VALUES ('7a4844ec-a616-416c-8cc0-c4607342cfad', 'reinhardluvr69', '2022-12-28 14:19:28.000000', null)
	_, err = t.db.Pool.Exec(ctx, `INSERT INTO "Token" (access_token, username, created_at, expires_at) 
	VALUES ($1, $2, $3, $4)`,
		token.AccessToken, token.Username, dt.Format(time.RFC3339), expiryTime.Format(time.RFC3339)) // todo: save
	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
	}
	// todo: scan
	return &token, err
}

func (t *tokensPostgresImpl) GetToken(ctx context.Context, accessToken string) (*turnip.Token, error) {
	token := turnip.Token{}
	var createdAt, expiresAt pgtype.Timestamp
	row := t.db.Pool.QueryRow(ctx, `SELECT *
FROM "Token" t
WHERE access_token=$1
LIMIT 1`, accessToken) // todo  get accessToken by access accessToken
	err := row.Scan(&token.AccessToken, &token.Username, createdAt, expiresAt)
	if err == pgx.ErrNoRows {
		// unauthorized or no token
		return nil, nil
	}
	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
	}

	if createdAt.Valid {
		token.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if expiresAt.Valid {
		token.ExpiresAt = timestamppb.New(expiresAt.Time)
	}

	return &token, nil
}

func NewTokensPostgres(ctx context.Context, d *clients.PostgresDb) Tokens {
	t := tokensPostgresImpl{
		db:          d,
		dbTableName: "Token",
	}

	clients.SetupTable(ctx, d, &t)

	return &t
}
