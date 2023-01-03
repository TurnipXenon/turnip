package migration

const (
	MigrateTags0001 = `create table "Tag"
(
    tag        text,
    content_id uuid
        constraint "Tag_Content_primary_id_fk"
            references "Content" ON DELETE CASCADE,
    created_at timestamp,
    constraint "Tag_pk"
        primary key (tag, content_id)
);

create index "Tag_content_id_index"
    on "Tag" (content_id);

create index "Tag_tag_index"
    on "Tag" (tag);`
)
