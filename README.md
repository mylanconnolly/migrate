# Migrate

This is a simple and opinionated PostgreSQL migration tool.

This came about because I was duplicating migration functionality across
multiple projects, and they all had the same behavior.

This tool uses `github.com/jackc/pgx`, so if you're using it, it won't pull in
unnecessary dependencies. On the flip side, if you're not, then it will (hence
the opinionated aspect).

## How it works

The user passes in a directory that contains SQL migration files. These files
are named in the following format:

```
$version_$name.up.sql
$version_$name.down.sql
```

The `$version` part defines the version number of the migration and `$name`
is used to define the name of the migration. This name is discarded, but should
be used to keep track of what the migrations are supposed to do.

If a migration ends in `up.sql`, it defines what should happen when you migrate
it up. If a migration ends in `down.sql`, it defines what should happen when you
migrate it down.

The tool will maintain a table called `schema_migrations` which has two columns:

| Column        | Description                                     |
| :------------ | :---------------------------------------------- |
| `version`     | The version number, extracted from the filename |
| `migrated_at` | The timestamp of when the migration occurred.   |

Only the most recent migration is shown in the table, in order to keep the table
as small and fast as possible.

All pending migrations are performed in a transaction so that if any migration
fails, they are rolled back without any changes.
