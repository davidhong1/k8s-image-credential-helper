APP_NAME=k8simagecredentialhelper
APP_BINARY=bin/$(APP_NAME)

all: build

test:
	go test -v ./...

dep:
	go mod download

.PHONY: build
build: ## build
	go build -o $(APP_BINARY) -v main.go

clean:
	go mod tidy

run:
	./run.sh $(APP_BINARY)