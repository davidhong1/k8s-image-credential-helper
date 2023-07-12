APP_NAME=k8simagecredentialhelper
BASE_PATH=../
APP_BINARY=$(BASE_PATH)bin/$(APP_NAME)


all: build

test:
	go test -v ./$(BASE_PATH)...

dep:
	cd $(BASE_PATH) && go mod download

.PHONY: build
build: ## build
	go build -o $(APP_BINARY) -v $(BASE_PATH)main.go

run:
	./run.sh $(APP_BINARY)