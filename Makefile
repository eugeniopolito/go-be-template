DB_URL=postgresql://root:secret@localhost:5432/app?sslmode=disable

create-new-migration:
	migrate create -ext sql -dir db/migration -seq $(SEQ_NAME)

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

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
	docker-compose up -d

docker-down:
	docker-compose down

db-docs:
	dbdocs build docs/db.dbml

db-schema:
	dbml2sql --postgres -o docs/schema.sql docs/db.dbml

create-swagger-doc:
	swag init

.PHONY: create-new-migration migrateup migratedown sqlc test server build build-docker-image docker-up docker-down db-docs db-schema create-swagger-doc