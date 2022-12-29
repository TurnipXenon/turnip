package migration

const (
	MigrateContent0001 = `create table "Content"
(
    primary_id     uuid not null
        constraint "Content_pk"
            primary key,
    created_at     timestamp,
    title          text,
    description    text,
    content        text,
    tag_list       text[],
    access_details text,
    meta           text,
    author_id      uuid
        constraint "Content_User_primary_id_fk"
            references "User"
);

alter table "Content"
    owner to postgres;`
)
