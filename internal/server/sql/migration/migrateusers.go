package migration

const (
	MigrateUsers0001 = `create table "User"
(
    primary_id      uuid not null
        constraint "User_pk"
            primary key,
    username        varchar(50)
        constraint "User_pk2"
            unique,
    hashed_password varchar(60),
    access_groups   text[]
);

alter table "User"
    owner to postgres;

create unique index "User_username_index"
    on "User" (username);`
)
