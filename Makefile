DB_URL=postgresql://root:secret@localhost:5432/app?sslmode=disable

init:
	go mod init github.com/eugeniopolito/gobetemplate

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root app

dropdb:
	docker exec -it postgres12 dropdb app

create-new-migration:
	migrate create -ext sql -dir db/migration -seq $(SEQ_NAME)

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc --experimental generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

build:
	go build -o app

build-docker-image:
	docker-compose build --no-cache --pull 

docker-up:
	docker-compose up

docker-down:
	docker-compose down

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

db-docs:
	dbdocs build docs/db.dbml

db-schema:
	dbml2sql --postgres -o docs/schema.sql docs/db.dbml

create-swagger-doc:
	swag init

.PHONY: create-docker-network createdb dropdb migrateup migratedown sqlc test server build migrateup1 migratedown1 build-docker-image docker-up docker-down redis db-docs db-schema create-swagger-doc