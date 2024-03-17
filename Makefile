.PHONY: build test clean

APP_PATH := ./cmd/main.go
APP_NAME := filmhub
TEST_PATH := ./internal/handler

build:
	docker-compose build $(APP_NAME)

run:
	docker-compose up $(APP_NAME)

test:
	go test -cover -v $(TEST_PATH)

clean:
	rm -f $(APP_NAME)
