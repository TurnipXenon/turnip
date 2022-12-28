// storage is an abstraction to s3 buckets

package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/TurnipXenon/turnip/internal/clients"
	"github.com/TurnipXenon/turnip/internal/server/sql/migration"
	"github.com/TurnipXenon/turnip/internal/util"
)

const (
	ifUserTableExists = `SELECT EXISTS(
               SELECT
               FROM information_schema.tables
               WHERE table_schema = 'public'
                 AND table_name = 'User'
           );`
)

type usersPostgresImpl struct {
	db           *clients.PostgresDb
	ddbTableName *string
}

func (u *usersPostgresImpl) CreateUser(ctx context.Context, ud *User) error {
	// check if user already exists
	rows := u.db.Pool.QueryRow(
		ctx,
		`SELECT EXISTS(
    			SELECT username FROM public."User" WHERE username=$1
           )`,
		ud.Username,
	)
	var exists bool
	err := rows.Scan(&exists)
	if err != nil {
		util.LogDetailedError(err)
		return util.WrapErrorWithDetails(err)
	}
	if exists {
		return util.WrapErrorWithDetails(UserAlreadyExists)
	}

	_, err = u.db.Pool.Exec(
		ctx,
		`INSERT INTO public."User" 
    			(primary_id, username, hashed_password, access_groups) 
				VALUES ($1, $2, $3, '{}')`,
		uuid.New().String(), // todo generate uuid
		ud.Username,
		ud.HashedPassword,
	) // todo GET User
	if err != nil {
		util.LogDetailedError(err)
		return util.WrapErrorWithDetails(err)
	}

	// todo: test this out

	return nil
}

func (u *usersPostgresImpl) GetUser(ctx context.Context, s *User) (*User, error) {
	// todo: figure out access group list
	row := u.db.Pool.QueryRow(
		ctx,
		`SELECT username FROM public."User" WHERE username=$1`,
		s.Username,
	)
	newUser := User{}
	err := row.Scan(&newUser.Username)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
	}

	return &newUser, nil
}

func NewUsersPostgres(ctx context.Context, d *clients.PostgresDb) Users {
	s := usersPostgresImpl{
		db:           d,
		ddbTableName: aws.String("Users"),
	}

	rows := s.db.Pool.QueryRow(ctx, ifUserTableExists)
	var exists bool
	err := rows.Scan(&exists)
	if err != nil {
		util.LogDetailedError(err)
		log.Fatalf("failed to check if table exists: %v", err)
	}
	if !exists {
		// from RocketDonkey @ https://stackoverflow.com/a/14668907/17836168
		s.migrateUsers(ctx)
	}

	// todo: detect schema change

	return &s
}

func (u *usersPostgresImpl) migrateUsers(ctx context.Context) {
	_, err := u.db.Pool.Exec(ctx, migration.MigrateUsers0001)
	if err != nil {
		util.LogDetailedError(err)
		log.Fatal(err)
	}
}
