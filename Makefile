DB_URL=postgres://requests:requests@postgres:5432/requests?sslmode=disable
TEST_DB_URL=postgres://requests:requests@postgres:5432/requests_test?sslmode=disable

up:
	docker compose up -d postgres

down:
	docker compose down

migrate-up:
	docker compose run --rm migrate -path=/migrations -database "$(DB_URL)" up

migrate-down:
	docker compose run --rm migrate -path=/migrations -database "$(DB_URL)" down -all

test-db:
	docker compose exec postgres psql -U requests -d postgres -c "SELECT 1 FROM pg_database WHERE datname = 'requests_test'" | grep -q 1 || docker compose exec postgres psql -U requests -d postgres -c "CREATE DATABASE requests_test"
	docker compose run --rm migrate -path=/migrations -database "$(TEST_DB_URL)" up

test: test-db
	TEST_DATABASE_URL=postgres://requests:requests@localhost:5432/requests_test?sslmode=disable go test ./...

lint:
	go vet ./...
