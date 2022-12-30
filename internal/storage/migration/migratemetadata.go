package migration

const (
	MigrateMetadata0001 = `create table "Metadata"
        (
            key   text
                constraint "Metadata_pk"
                    primary key,
            value text
        )`
)
