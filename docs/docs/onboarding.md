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
    go run cmd/turnip.go
    ```

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
3. If you edited index.md, run `go run scripts/sync_mkdocs_readme.go` from the root folder `/`.
