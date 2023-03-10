package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/twitchtv/twirp"

	"github.com/TurnipXenon/turnip/internal/config"
	"github.com/TurnipXenon/turnip/internal/storage/migration"
	"github.com/TurnipXenon/turnip/internal/util"
)

const (
	// todo:  note!!!
	usernameCharLimit = 50
)

type usersPostgresImpl struct {
	db          *PostgresDb
	dbTableName string
	sysConf     config.SystemConfig
}

func (u *usersPostgresImpl) GetTableName() string {
	return u.dbTableName
}

func (u *usersPostgresImpl) GetMigrationSequence() []migration.Migration {
	return []migration.Migration{
		migration.NewGenericMigration(migration.MigrateUsers0001),
	}
}

func (u *usersPostgresImpl) CreateUser(ctx context.Context, ud *User) error {
	if !u.sysConf.CanUserBeMade.Get() {
		return twirp.Unavailable.Error("CreateUser endpoint is unavailable")
	}

	// check if user already exists
	rows := u.db.Pool.QueryRow(
		ctx,
		`SELECT EXISTS(
    			SELECT username FROM "User" WHERE username=$1
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
		`INSERT INTO "User" 
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
		`SELECT username, hashed_password, primary_id FROM "User" WHERE username=$1`,
		s.Username,
	)
	newUser := User{}
	err := row.Scan(&newUser.Username, &newUser.HashedPassword, &newUser.PrimaryId)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		util.LogDetailedError(err)
		return nil, util.WrapErrorWithDetails(err)
	}

	return &newUser, nil
}

func NewUsersPostgres(ctx context.Context, d *PostgresDb, conf config.SystemConfig) Users {
	p := usersPostgresImpl{
		db:          d,
		dbTableName: "User",
		sysConf:     conf,
	}

	SetupTable(ctx, d, &p)

	// todo(turnip): detect schema change

	return &p
}
