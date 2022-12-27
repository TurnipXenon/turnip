package migration

const (
	MigrateUsers0001 = `create table public."User"
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

alter table public."User"
    owner to postgres;

create unique index "User_username_index"
    on public."User" (username);`
)
