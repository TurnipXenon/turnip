# Onboarding

## Setup

**Required:**

- Go 1.17
- Docker

**Steps:**

1. Run docker.
    ```shell
    docker-compose up
    ```
    - This command runs all the dependencies that our Go program will call.
2. Install the go dependencies
    ```shell
    go install
    ```
3. Run the go program
    ```shell
    go run cmd/turnip.go -is-local
    ```

Locally, these three ports have special meanings:

- 8000
- 8010
- 8020

Using any other port may not work.

If no port argument was given, we default to port 8000.

Envs:
```
PGCONN
DATABASE_URL by DigitalOcean
CORS_ALLOWLIST=["", ""]
POTATO_TOKEN
POTATO_URL
```

**todo: postgres onboarding**

logs:

```
postgres=# CREATE DATABASE turnip
postgres-# \connect turnip
connection to server at "localhost" (::1), port 5432 failed: FATAL:  database "turnip" does not exist
Previous connection kept
postgres-# CREATE role turnipservice WITH LOGIN PASSWORD 'password';
ERROR:  syntax error at or near "CREATE"
LINE 2: CREATE role turnipservice WITH LOGIN PASSWORD 'password';
^
postgres=# CREATE ROLE turnipservice WITH LOGIN PASSWORD 'password';
CREATE ROLE
postgres=# ALTER ROLE turnipservice CREATEDB;
ALTER ROLE
postgres=# CREATE SCHEMA public;
ERROR:  schema "public" already exists
postgres=# GRANT ALL ON SCHEMA public to turnipservice;
GRANT
postgres=# GRANT ALL ON SCHEMA public to public;
GRANT
postgres=# \q
Press any key to
```

Need this one too: `GRANT postgres to turnipservice;`

YOU GOTTA SHARE THE FACT YOUVE BEEN TROUBLESHOOTING FOR TWO HOURS when DigitalOcean was apparently giving you the wrong
IP address. or maybe something is up with my browser. then again, some rando website https://whatismyipaddress.com/ found my ip correctly.
DBMS: PostgreSQL (ver. 12.0)Case sensitivity: plain=mixed, delimited=exact The connection attempt failed.


## MKDocs

You don't really need to run through this to make edits to MKDocs, but if you want to see the layout and what it looks
like served, check this guide.

This assumes that you have **Python** installed locally.

### MKDocs: Setup

```shell
pip install mkdocs
```

### MKDocs: Commands

When entering these commands, go to `/docs` instead of being in the project's root folder `/`.

* `mkdocs serve` - Start the live-reloading docs server.
* `mkdocs build` - Build the documentation site.
* `mkdocs -h` - Print help message and exit.

### MKDocs: Ideal workflow

1. Make changes
2. See changes made using `mkdocs serve`
3. If you edited index.md, run `go run dev/sync_mkdocs_readme.go` from the root folder `/`.
