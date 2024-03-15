.PHONY: build test clean

APP_PATH := ./cmd/main.go
APP_NAME := filmhub

build:
	go build -o $(APP_NAME) $(APP_PATH)

run: build
	go run $(APP_PATH)

test:
	go test -cover -v ./internal/handler

clean:
	rm -f $(APP_NAME)
