.PHONY: all build run test clean docker-build docker-run fmt vet lint

APP_NAME=taiwanstay-back
DOCKER_IMAGE=$(APP_NAME)

all: build

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

test:
	go test -v ./...

clean:
	rm -rf bin

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	# Requires golangci-lint installed
	golangci-lint run ./...

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)
