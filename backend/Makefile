# Variables
APP_NAME = stock-information
DOCKER_IMAGE = stock-information
DB_NAME = stockdb
COCKROACH_CONTAINER = cockroachdb
COCKROACH_PORT = 26257
COVERAGE_FILE = coverage.out
COVERAGE_REPORT = coverage.html

.PHONY: all build run docker-build docker-up db-create clean stop

build:
	go build -o $(APP_NAME) ./cmd/server

run: build
	./$(APP_NAME)

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-up:
	docker compose up --build

db-create:
	docker exec -it $$(docker ps -qf "name=$(COCKROACH_CONTAINER)") \
	./cockroach sql --insecure --execute="CREATE DATABASE IF NOT EXISTS $(DB_NAME);"

clean:
	docker compose down -v --remove-orphans
	docker image rm $(DOCKER_IMAGE) || true

stop:
	docker compose down

dev:
	air

format:
	gofmt -w .

test:
	go test -coverprofile=$(COVERAGE_FILE) ./...

coverage: test
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_REPORT)
