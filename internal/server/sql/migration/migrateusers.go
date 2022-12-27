package migration

const (
	MigrateUsers0001 = `create table "User"
(
    primary_id      uuid not null
        constraint "User_pk"
            primary key,
    username        char(50)
        constraint "User_pk2"
            unique,
    hashed_password char(120),
    access_groups   text[]
);

create unique index "User_username_index"
    on "User" (username);`
)
