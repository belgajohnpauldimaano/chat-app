migrateup:
	migrate -path infrastructure/db/migrations -database "postgresql://postgres:postgres@localhost:5433/chat_app_db?sslmode=disable" -verbose up

migratedown:
	migrate -path infrastructure/db/migrations -database "postgresql://postgres:postgres@localhost:5433/chat_app_db?sslmode=disable" -verbose down

migrateup-docker:
	docker run --rm -v ./infrastructure/db/migrations:/migrations --network host migrate/migrate -path=/migrations -database "postgres://postgres:postgres@localhost:5433/chat_app_db?sslmode=disable" -verbose up

migratedown-docker:
	docker run --rm -v ./infrastructure/db/migrations:/migrations --network host migrate/migrate -path=/migrations -database "postgres://postgres:postgres@localhost:5433/chat_app_db?sslmode=disable" -verbose down

create-migration-file:
	echo "Creating migration file..."
	docker run --rm -v ./infrastructure/db/migrations:/migrations migrate/migrate create -ext sql -dir /migrations $(migration_name)

run-app:
	go run cmd/main.go

run-air:
	air

.PHONY: migrateup migratedown run-app run-air migrateup-docker migratedown-docker create-migration-file