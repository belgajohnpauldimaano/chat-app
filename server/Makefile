migrateup:
	migrate -path infrastructure/db/migrations -database "postgresql://postgres:postgres@localhost:5433/chat_app_db?sslmode=disable" -verbose up

migratedown:
	migrate -path infrastructure/db/migrations -database "postgresql://postgres:postgres@localhost:5433/chat_app_db?sslmode=disable" -verbose down

run-app:
	go run cmd/main.go

run-air:
	air

.PHONY: migrateup migratedown run-app run-air