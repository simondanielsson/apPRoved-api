.PHONY: build build-linux run build-image deploy test docs
IMAGE_NAME := approved-api
IMAGE_TAG := latest

build:
	@go build -o bin/approved cmd/main.go

build-linux:
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/approved cmd/main.go

docs:
	swag init -g cmd/main.go

run: build
	@./bin/approved

run-dev: docs build
	@go run cmd/main.go

build-image: docs
	echo "Building ${IMAGE_NAME}:${IMAGE_TAG} image...";
	docker build -f arm.Dockerfile -t ${IMAGE_NAME}:${IMAGE_TAG} .;
	echo "Image ${IMAGE_NAME}:${IMAGE_TAG} built.";

deploy:
	echo "Running docker compose"
	docker compose up -d

test:
	@go test -v ./...
