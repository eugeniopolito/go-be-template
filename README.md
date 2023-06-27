# Golang Back End Template

A Docker multi container for APIs development in Go!
This framework includes everything you need to build APIs with this amazing language: PostgreSQL, Redis and Go Gin framework for rapid Back End development!

## Setup local development

### Install tools

- [Docker desktop](https://www.docker.com/products/docker-desktop)
- [TablePlus](https://tableplus.com/)
- [Golang](https://golang.org/)
- [Homebrew](https://brew.sh/)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

    ```bash
    brew install golang-migrate
    ```

- [Sqlc](https://github.com/kyleconroy/sqlc#installation)

    ```bash
    brew install sqlc
    ```

- [DB Docs](https://dbdocs.io/docs)

    ```bash
    npm install -g dbdocs
    dbdocs login
    ```

- [DBML CLI](https://www.dbml.org/cli/#installation)

    ```bash
    npm install -g @dbml/cli
    dbml2sql --version
    ```


### Setup infrastructure

- Copy the ```app.env.example``` to ```app.env``` and fill in your SMTP and DB preferences.

- Start postgres container:

    ```bash
    make postgres
    ```

- Create database:

    ```bash
    make createdb
    ```

- Create Redis:

    ```bash
    make redis
    ```

- Create new migration:

    ```bash
    migrate create -ext sql -dir db/migration -seq $(SEQ_NAME)
    ```

- Run db migration up all versions:

    ```bash
    make migrateup
    ```

- Run db migration up 1 version:

    ```bash
    make migrateup1
    ```

- Run db migration down all versions:

    ```bash
    make migratedown
    ```

- Run db migration down 1 version:

    ```bash
    make migratedown1
    ```

### Documentation

- Generate DBML documentation:

    ```bash
    make db-docs
    ```

### How to generate code

- Generate schema SQL file with DBML:

    ```bash
    make db-schema
    ```

- Generate SQL CRUD with sqlc:

    ```bash
    make sqlc
    ```

- Create a new db migration:

    ```bash
    make create-new-migration SEQ_NAME=<migration_name>
    ```

- Create a new DBML documentation:

    ```bash
    make db_docs
    ```

- Create a new schema from DBML:

    ```bash
    make db_schema
    ```

### How to run

- Run server:

    ```bash
    make server
    ```

- Run test:

    ```bash
    make test
    ```

- Run Dockerized app:

    ```bash
    make docker-up
    ```

- Visit the Swagger documentation at local address:

    ```
    http://localhost:8080/swagger/index.html
    ```
