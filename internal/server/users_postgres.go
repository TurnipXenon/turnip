// storage is an abstraction to s3 buckets

package server

import (
	"context"
	"errors"
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
	// check if migration already exists
	row, err := u.db.Conn.Query(ctx, "INSERT") // todo GET User
	row = row                                  // prevent compile error
	if err != nil {
		util.LogDetailedError(err)
		return util.WrapErrorWithDetails(err)
	}
	// check if exists
	// todo row

	_, err = u.db.Conn.Exec(ctx, "") // todo PUT migration
	if err != nil {
		util.LogDetailedError(err)
		return util.WrapErrorWithDetails(err)
	}

	return nil

	//TODO implement me
	panic("implement me")
}

func (u *usersPostgresImpl) GetUser(s *User) (*User, error) {
	//TODO implement me
	panic("implement me")
}

func NewUsersPostgres(ctx context.Context, d *clients.PostgresDb) Users {
	s := usersPostgresImpl{
		db:           d,
		ddbTableName: aws.String("Users"),
	}

	// todo: check if table exists
	rows, err := s.db.Conn.Query(ctx, ifUserTableExists)
	if err != nil {
		// todo
		util.LogDetailedError(errors.New("TODO: handle error for table search"))
		log.Fatalf("TODO: handle error for table search")
	}
	defer rows.Close()

	exists := false
	for rows.Next() {
		err = rows.Scan(&exists)
	}
	if err != nil || !exists {
		// from RocketDonkey @ https://stackoverflow.com/a/14668907/17836168
		s.migrateUsers(ctx)
	}

	// todo: detect schema change

	return &s
}

func (u *usersPostgresImpl) migrateUsers(ctx context.Context) {
	_, err := u.db.Conn.Exec(ctx, migration.MigrateUsers0001)
	if err != nil {
		util.LogDetailedError(err)
		log.Fatal(err)
	}
}
