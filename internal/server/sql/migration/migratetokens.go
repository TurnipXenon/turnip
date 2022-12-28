package migration

const (
	MigrateToken0001 = `create table "Token"
(
    access_token uuid        not null
        constraint "Token_pk"
            primary key,
    username     varchar(50) not null
        constraint "Token_User_username_fk"
            references "User" (username),
    created_at   timestamp,
    expires_at   timestamp
);

create index "Token_username_index"
    on "Token" (username);

`
)
