dockerup:
	docker compose up

dockerdown:
	docker compose down

createdb:
	docker exec -it postgres_db createdb --username=admin --owner=admin simple_bank

dropdb:
	docker exec -it postgres_db dropdb --username=admin simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://admin:password@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://admin:password@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

start:
	go run main.go