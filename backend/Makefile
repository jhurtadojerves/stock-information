# Variables
DOCKER_IMAGE = stock-information
DB_NAME = stockdb
COVERAGE_FILE = coverage.out
COVERAGE_REPORT = coverage.html

.PHONY: all build run docker-build docker-up db-create clean stop

build:
	docker build -t $(DOCKER_IMAGE) .

up:
	docker compose up --build

clean:
	docker compose down -v --remove-orphans
	docker image rm $(DOCKER_IMAGE) || true

stop:
	docker compose down

test:
	go test -coverprofile=$(COVERAGE_FILE) ./...

coverage: test
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_REPORT)
