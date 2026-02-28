.PHONY: build run test docker-up docker-down migrate generate lint clean

# Build the application
build:
	go build -o bin/api ./cmd/api

# Run the application
run:
	go run ./cmd/api

# Run tests
test:
	go test -v ./...

# Start Docker containers
docker-up:
	docker-compose up -d

# Stop Docker containers
docker-down:
	docker-compose down

# Run database migrations
migrate:
	@echo "Running migrations..."
	@for f in db/migrations/*.sql; do \
		echo "Executing $$f"; \
		PGPASSWORD=$(DB_PASS) psql -h $(DB_HOST) -U $(DB_USER) -d $(DB_NAME) -f $$f; \
	done

# Generate API code from OpenAPI spec
generate:
	oapi-codegen -package generated -generate types,server,spec api/openapi.yaml > internal/api/generated/models.gen.go

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/
