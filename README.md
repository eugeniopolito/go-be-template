# Golang Back End Template

A Docker multi container for APIs development in Go!
This framework includes everything you need to build APIs with Golang: Nginx, PostgreSQL, Redis and Gin framework for rapid Back End development!

## Setup local development

### Install tools

- [Docker desktop](https://www.docker.com/products/docker-desktop)
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

- Build the docker image:

    ```bash
    make build-docker-image
    ```

- Start all required Docker containers:

    ```bash
    make docker-up
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
    make db-docs
    ```

- Create a new schema from DBML:

    ```bash
    make db-schema
    ```

- Create Swagger documentation:

    ```bash
    make create-swagger-doc
    ```

### How to run

- Run test:

    ```bash
    make test
    ```

- Run the ```main.go``` as standalone server:

    ```bash
    make server
    ```

- Or run it as Dockerized app with Docker Compose:

    ```bash
    make docker-up
    ```

- Visit the Swagger APIs documentation at your local address:

    ```
    http://localhost/swagger/index.html
    ```
